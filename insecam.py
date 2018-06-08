#!/usr/bin/env python3
from lxml import html
import requests
import pickle

entry_url = 'http://www.insecam.org'
headers = {'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.96 Safari/537.36'}

cams_model = {'cams': [], 'page_counter': 0}

def getpage(tz, index):
    page = requests.get(entry_url + '/en/bytimezone/' + tz + '/?page=' + str(index), headers=headers)
    if page.status_code != 200:
        return None, False
    tree = html.fromstring(page.content)
    cams = tree.xpath('//img[@class="thumbnail-item__img img-responsive"]/@src')
    cams = list(map(lambda x: (x, tz), cams))
    return cams, True

def save_cam(cams):
    pickle.dump(cams, open('insecam_tz.p', 'wb'))

if __name__ == '__main__':
    cams = []
    tz_list = ['+00:00',  '+01:00', '+02:00', '+03:00', '+03:30', '+04:00', '+04:30', '+05:00', '+05:30', '+05:45', '+06:00', '+07:00', '+08:00', '+09:00', '+09:30', '+10:00', '+10:30', '+11:00', '+12:00', '+13:00', '-', '-02:00', '-03:00', '-03:30', '-04:00', '-04:30', '-05:00', '-06:00', '-07:00', '-08:00', '-09:00', '-10:00']
    for tz in tz_list:
        i = 0
        ok = True
        last_cams = None
        while ok:
            lst, ok = getpage(tz, i)
            print('TZ: %s, page %d (total cam: %d)' % (tz, i, len(cams)))
            if last_cams == lst:
                break
            last_cams = lst
            cams += lst
            i += 1
    save_cam(cams)
