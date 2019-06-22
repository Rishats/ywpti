# Weather Notifications

Simple system which notify users about weather in to Telegram via [Horn](https://github.com/requilence/integram)

### Installing
```
 cp .env.example .env
```
And configure your .env
```YW_API_KEY=be6653b5-4fdd-41a9-a31c-b3a935252493
   YW_API_URI=https://api.weather.yandex.ru/v1/informers
   YW_LAT=43.262547
   YW_LON=76.927127
   YW_LANG=kk_KZ
   INTEGRAM_WEBHOOK_URI=https://integram.org/webhook/cCgds28sIpG
   ```

Run via docker,
```
 docker build -t wn .
 docker run -d wn
```

Or via go native

```
go run main.go
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/Rishats/ywpti/tags). 

## Authors

* **Rishat Sultanov** - [Rishats](https://github.com/Rishats)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
