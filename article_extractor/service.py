#!/usr/bin/env python
# encoding: utf-8

import simplejson as json
from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
import requests

from article_extractor import ArticleExtractor

article_extractor = None

PORT = 1030


class ArticleService(BaseHTTPRequestHandler):

    def do_POST(self):
        content_len = int(self.headers.getheader('content-length'))
        raw_text = self.rfile.read(content_len)
        raw_text = raw_text.decode("utf8")
        article_result = self.extract(raw_text)

        self.send_response(200)
        self.send_header("Content-Type", "text/javascript; charset=UTF-8")
        self.end_headers()

        self.wfile.write(json.dumps(article_result))
        return

    def extract(self, url):
        return article_extractor.extract(url)


def article_extract(url):
    r = requests.post("http://localhost:%d/" % (PORT), data=url,
                      headers={'content-type': 'text/plain; chartset=utf-8'})
    if r.status_code != requests.codes.ok:
        r.raise_for_status()

    return json.loads(r.text)


def main():
    # only setup extractors if we're running the server
    global article_extractor
    article_extractor = ArticleExtractor()

    server = HTTPServer(('', PORT), ArticleService)
    server.serve_forever()

if __name__ == "__main__":
    main()
