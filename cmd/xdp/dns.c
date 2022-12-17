#include <linux/if_ether.h>
#include <linux/in.h>
#include <bcc/proto.h>
int dns_matching(struct __sk_buff *skb) {
    u8 *cursor = 0;
     // Checking the IP protocol:
    struct ethernet_t *ethernet = cursor_advance(cursor, sizeof(*ethernet));
    if (ethernet->type == ETH_P_IP) {
         // Checking the UDP protocol:
        struct ip_t *ip = cursor_advance(cursor, sizeof(*ip));
        if (ip->nextp == IPPROTO_UDP) {
             // Check the port 53:
            struct udp_t *udp = cursor_advance(cursor, sizeof(*udp));
            if (udp->dport == 53 || udp->sport == 53) {
                return -1;
            }
        }
    }
    return 0;
}