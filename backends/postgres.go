package backends

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	// pgsql
	_ "github.com/lib/pq"
	"github.com/millken/mkdns/types"
)

func init() {
	Register("postgres", NewPostgresBackend)
}

// PostgresBackend postgres as backend store
type PostgresBackend struct {
	db        *sql.DB
	url       *url.URL
	lastUtime string
}

// NewPostgresBackend init
func NewPostgresBackend(u *url.URL) (Backend, error) {
	db, err := sql.Open("postgres", u.String())
	if err != nil {
		return nil, fmt.Errorf("[%s]: %s", u.String(), err)
	}
	backend := &PostgresBackend{
		db:  db,
		url: u,
	}
	return backend, nil
}

func (b *PostgresBackend) read(path string) {
}

func (b *PostgresBackend) getMaxCount() (maxID, countID int, err error) {
	err = b.db.QueryRow("select max(id) max_id,count(id) count_id from config.record").Scan(&maxID, &countID)
	return
}

func (b *PostgresBackend) getMaxTime() (utime string, err error) {
	err = b.db.QueryRow("select max(utime) utime from config.event").Scan(&utime)
	return
}

// Load load config
func (b *PostgresBackend) Load() {

	var id int
	var domain string
	var value []byte
	log.Printf("[INFO] loading postgres : %s", b.url.String())
	utime, err := b.getMaxTime()
	if err != nil {
		log.Printf("[ERROR] getMaxTime() : %s", err)
		return
	}
	b.lastUtime = utime
	maxID, countID, err := b.getMaxCount()
	if err != nil {
		log.Printf("[ERROR] getMaxTime() : %s", err)
		return
	}

	log.Printf("[INFO] postgres utime=%s, maxID=%d, countID=%d", utime, maxID, countID)
	offset := 0
	limit := 2500
	for offset < countID {
		sqlstr := fmt.Sprintf("SELECT * FROM config.record where id<=$1 order by id asc limit %d offset %d", limit, offset)
		rows, err := b.db.Query(sqlstr, maxID)
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
		defer rows.Close()
		for rows.Next() {

			err = rows.Scan(&id, &domain, &value)
			if err != nil {
				log.Fatal(err)
			}
			offset++
			dpb, err := types.DecodeByProtobuff(value)
			if err != nil {
				log.Printf("[ERROR] DecodeByProtobuf err: %s", err)
				continue
			}

			if dpb.Domain != "" {
				domain = dpb.Domain
			}
			zonemap.Set(domain, value)
			zonecache.Del(domain)
			log.Printf("[DEBUG] ID: %d, domain: %s, dpb=%+v\n", id, domain, dpb)
		}
	}
	go b.watch()

}

func (b *PostgresBackend) watch() {
	var domain, utime, act string
	var value []byte
	var isLocked bool
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if isLocked {
				continue
			}
			sqlstr := "select event.domain, utime, act, value from config.event left outer join config.record on event.domain = record.domain where utime>$1 order by utime asc"
			rows, err := b.db.Query(sqlstr, b.lastUtime)
			if err != nil {
				log.Printf("[ERROR] watch query %s", err)
			}
			isLocked = true
			for rows.Next() {
				err = rows.Scan(&domain, &utime, &act, &value)
				if err != nil {
					log.Printf("[ERROR] rows scan %s", err)
				}
				if act == "delete" {
					log.Printf("[INFO] domain config removed : %s", domain)
					zonemap.Remove(domain)
					zonecache.Del(domain)
					b.lastUtime = utime
					continue
				}
				dpb, err := types.DecodeByProtobuff(value)
				if err != nil {
					log.Printf("[ERROR] DecodeByProtobuf err: %s", err)
					continue
				}

				if dpb.Domain != "" {
					domain = dpb.Domain
				}
				zonemap.Set(domain, value)
				zonecache.Del(domain)
				b.lastUtime = utime
				log.Printf("[INFO] domain config added : %s", domain)
			}
			isLocked = false
		}
	}
}
