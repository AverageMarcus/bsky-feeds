# jetstream-feed-generator

Glues together [Jetstream](https://github.com/bluesky-social/jetstream) and [go-bsky-feed-generator](https://github.com/ericvolp12/go-bsky-feed-generator/) with some SQLite to consume the Bluesky firehose and serve a feed based on posts matching some criteria.

This was forked from [rolandcrosby/jetstream-feed-generator](https://github.com/rolandcrosby/jetstream-feed-generator) (huge thank you to rolandcrosby!) and expanded with the following features:
* Supports running multiple feeds
* Added support for using PostgreSQL as well as SQLite
* Handled retries on some error conditions
* Dynamically create feeds based on the provided config
* Implemented a generic "SimpleFeed" that can dynamically create feeds from config

At time of writing, this is currently being used to power the following feeds:
* [KubeCon & Related Events](https://bsky.app/profile/did:plc:mtepw4cvbmdvu7zygmm5xbop/feed/kubecon)
* [KubeCon Parties](https://bsky.app/profile/did:plc:mtepw4cvbmdvu7zygmm5xbop/feed/kubecon-party)
* [EMF Camp](https://bsky.app/profile/did:plc:mtepw4cvbmdvu7zygmm5xbop/feed/emf-camp)

---

## Example Config

```yaml
db:
  engine: sqlite
  connection_string: feeds.sqlite
feed_names:
- composer-errors
- english-text
# If you weant to set up a feed that just watches for words/hashtags/accounts
# without any complex logic then use simple_feeds
simple_feeds:
  - name: emf-camp
    hashtags: ["EMFCamp", "EMF2026", "EMF26"]
    account_dids: ["did:plc:r5tbkz2suj4hz6kyadj73y6n"]
    regex: "(?mi)(^|\\s|#)(EMF ?Camp)(\\d{2,4})?(\\W|$)"
    include_year: true
log_level: INFO
log_format: json
consumer:
  enabled: true
  jetstream_url: wss://jetstream1.us-east.bsky.network/subscribe
  start_cursor: 0
feedgen:
  enabled: true
  port: 9072
  feed_actor_did: did:plc:replace-me-with-your-did
  service_endpoint: https://replace-me-with-your-service-endpoint.example.com
```

## Running

Run with:

```shell
go run . --config ./config.yml
```

The following URLs will then be available:

* http://localhost:9072/.well-known/did.json
* http://localhost:9072/xrpc/app.bsky.feed.describeFeedGenerator

Then for each feed enabled there will be an URL such as:
* http://localhost:9072/xrpc/app.bsky.feed.getFeedSkeleton?feed=at://did:plc:replace-me-with-your-did/app.bsky.feed.generator/english-text

## Adding a Feed to Your Profile

I recommend using the CLI from [bluesky-social/feed-generator](https://github.com/bluesky-social/feed-generator) to add the feed to your profile.


1. Checkout: <https://github.com/bluesky-social/feed-generator>
2. `yarn install`
3. Create a `.env` file with the following contents:

   ```bash
   FEEDGEN_HOSTNAME="The hostname you are hosting the generated feeds on"
   FEEDGEN_PUBLISHER_DID="Your accounts DID"
   FEEDGEN_SUBSCRIPTION_ENDPOINT="wss://bsky.network"
   FEEDGEN_SUBSCRIPTION_RECONNECT_DELAY=3000
   ```
4. `yarn publishFeed` and answer the questions

**Notes:**

* The “short name” is the name of the feed you specified in the config
* For the “avatar” you provide a path to a file on local disk.
