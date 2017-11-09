package plugins

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/abh/geoip"
)

type GeoIP struct {
	dbDirectory     string
	country         *geoip.GeoIP
	hasCountry      bool
	countryLastLoad time.Time

	city         *geoip.GeoIP
	cityLastLoad time.Time
	hasCity      bool

	asn         *geoip.GeoIP
	hasAsn      bool
	asnLastLoad time.Time
}

var geoIP = new(GeoIP)

func init() {
	//geoIP.setDbDirectory("/home/millken/qbox_sync/db")
	geoIP.setupGeoIPCountry()
}

func (g *GeoIP) GetCountry(ip net.IP) (country, continent string, netmask int) {
	if g.country == nil {
		return "", "", 0
	}

	country, netmask = geoIP.country.GetCountry(ip.String())
	if len(country) > 0 {
		country = strings.ToLower(country)
		continent = CountryContinent[country]
	}
	return
}

func (g *GeoIP) GetCountryRegion(ip net.IP) (country, continent, regionGroup, region string, netmask int) {
	if g.city == nil {
		log.Println("No city database available")
		country, continent, netmask = g.GetCountry(ip)
		return
	}

	record := geoIP.city.GetRecord(ip.String())
	if record == nil {
		return
	}

	country = record.CountryCode
	region = record.Region
	if len(country) > 0 {
		country = strings.ToLower(country)
		continent = CountryContinent[country]

		if len(region) > 0 {
			region = country + "-" + strings.ToLower(region)
			regionGroup = CountryRegionGroup(country, region)
		}

	}
	return
}

func (g *GeoIP) GetASN(ip net.IP) (asn string, netmask int) {
	if g.asn == nil {
		log.Println("No asn database available")
		return
	}
	name, netmask := g.asn.GetName(ip.String())
	if len(name) > 0 {
		index := strings.Index(name, " ")
		if index > 0 {
			asn = strings.ToLower(name[:index])
		}
	}
	return
}

func (g *GeoIP) setDbDirectory(dir string) {
	g.dbDirectory = dir
}

func (g *GeoIP) setDirectory() {
	if len(g.dbDirectory) > 0 {
		geoip.SetCustomDirectory(g.dbDirectory)
	}
}

func (g *GeoIP) setupGeoIPCountry() {
	if g.country != nil {
		return
	}

	g.setDirectory()

	gi, err := geoip.OpenType(geoip.GEOIP_COUNTRY_EDITION)
	if gi == nil || err != nil {
		log.Printf("Could not open country GeoIP database: %s\n", err)
		return
	}
	g.countryLastLoad = time.Now()
	g.hasCity = true
	g.country = gi

}

func (g *GeoIP) setupGeoIPCity() {
	if g.city != nil {
		return
	}

	g.setDirectory()

	gi, err := geoip.OpenType(geoip.GEOIP_CITY_EDITION_REV1)
	if gi == nil || err != nil {
		log.Printf("Could not open city GeoIP database: %s\n", err)
		return
	}
	g.cityLastLoad = time.Now()
	g.hasCity = true
	g.city = gi

}

func (g *GeoIP) setupGeoIPASN() {
	if g.asn != nil {
		return
	}

	g.setDirectory()

	gi, err := geoip.OpenType(geoip.GEOIP_ASNUM_EDITION)
	if gi == nil || err != nil {
		log.Printf("Could not open ASN GeoIP database: %s\n", err)
		return
	}
	g.asnLastLoad = time.Now()
	g.hasAsn = true
	g.asn = gi

}
