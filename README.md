# email-google-alerts

Looking at Google Alert emails:

1. `fetch` (data goes into `keywordalerts` mongo)
2. `gather` (data goes into `gatheredalerts` mongo)
3. `rank` (data goes into `rankedalerts` mongo)
4. `format` (data goes into csv)

Running

1. Run `article_extractor` (`nohup python service.py &`)
2. Run `fetch`
3. Run `gather`
4. Run `rank`
5. Run `format`
