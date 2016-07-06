# -*- coding: utf-8 -*-
import string
import re

# Third-party app imports
from urlparse import urlparse
from newspaper import Article


class ArticleExtractor(object):

    def __init__(self):
        self.title = ''

    def url_validate(self, url):
        url = urlparse(url)
        return (
            url.scheme + '://' + url.netloc +
            url.path, url.scheme + '://' + url.netloc
        )

    def extract(self, url):
        article = Article(url)
        article.download()
        article.parse()
        article.nlp()

        url, publisher = self.url_validate(url)

        data = {}
        data['url'] = url
        data['name'] = article.title  # Get Title
        if article.publish_date:
            data['created_at'] = str(article.publish_date)
        data['header_image'] = article.top_image
        data['basic_summary'] = article.summary
        data['opening_paragraph'] = article.opening_paragraph
        data['keywords'] = article.keywords
        data['authors'] = article.authors
        data['html'] = article.html
        data['text'] = article.text

        return data
