import pickle
import json
from urllib.parse import urlparse
import geoip2.database
import simplekml

kml=simplekml.Kml()

cams = pickle.load(open('insecam_tz.p', 'rb'))
reader = geoip2.database.Reader('GeoLite2-City_20170606.mmdb')

c_l = []
for cam in cams:
    c_l.append(cam[0])
    u = urlparse(cam[0])
    tz = cam[1]
    dsc = "Link: %s\nTZ: %s" % (cam[0], tz)
    host = u.netloc.split(':')[0]
    info = reader.city(host)
    kml.newpoint(name=host,
                 coords=[(info.location.longitude, info.location.latitude)],
                 description=dsc)
print(json.dumps(c_l))
kml.save('cams.kml')
