J
@¨SOA*=2;
ns1.test1.com.dns-admin.test1.com.¯¨—ë ÷ËH(Âé&0±&8’·+
@¨NS*
ns1.test1.com.
ns2.test2.com
@¨A*
	1.47.46.2
1.2.3.3
@¨TXT*
	AbCdFx11.0
@¨MX*
mx1.test.com.8*
mx2.test.com.8A
view¨A *
	1.47.46.2
1.2.3.3any*
1.2.3.4
1.2.4.5dxK
weight¨A *
	1.47.46.2
1.2.3.3*
1.2.3.4
1.2.4.5*
7.7.7.7
†
geo¨A *	
7.7.7.7*
1.7.6.2"asia*
1.7.6.5"asia*cn*
1.7.6.6*cn*
1.2.3.4
1.2.4.5*kr**
1.1.1.1
1.2.2.3
1.1.1.2"north-america*
1.1.1.3*us/
spf¨SPF* 
 v=spf1 include:spf.c3t.a -all2

_http._tcp¨SRV*
10 5 8080 www.example.com.