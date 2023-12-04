# Changelog

## 1.0.0 (2023-12-04)


### Features

* Init ([b0ee650](https://github.com/soerenschneider/fetcharr/commit/b0ee65047e18fa6310f0217bd6950bf79c26067e))

## [1.4.0](https://github.com/soerenschneider/fetcharr/compare/v1.3.1...v1.4.0) (2023-12-04)


### Features

* parse filenames of transferred files ([fcc373a](https://github.com/soerenschneider/fetcharr/commit/fcc373aa669cf6d6c8c1cecdc0a15d34ab7bdb8e))


### Bug Fixes

* **deps:** bump alpine from 3.18.4 to 3.18.5 ([1d74fbf](https://github.com/soerenschneider/fetcharr/commit/1d74fbfa6adf7a917048e5ac7e9ba271baba8666))
* **deps:** bump github.com/hashicorp/go-retryablehttp ([955c696](https://github.com/soerenschneider/fetcharr/commit/955c696ce9f78fdeef27a088d8f96d323953535e))
* **deps:** bump github.com/segmentio/kafka-go from 0.4.44 to 0.4.45 ([ec734f3](https://github.com/soerenschneider/fetcharr/commit/ec734f379b5bc01cfa1608b32fd6fec134ec0ab3))
* **deps:** bump github.com/segmentio/kafka-go from 0.4.45 to 0.4.46 ([969ed6c](https://github.com/soerenschneider/fetcharr/commit/969ed6c28898cbfb81e1dd5d319b3d12894c06ec))
* **deps:** bump golang from 1.21.3 to 1.21.4 ([53c1a50](https://github.com/soerenschneider/fetcharr/commit/53c1a500ae6fc2495aba19ceaf8f74411dfae639))

## [1.3.1](https://github.com/soerenschneider/fetcharr/compare/v1.3.0...v1.3.1) (2023-11-07)


### Bug Fixes

* **deps:** bump github.com/go-playground/validator/v10 ([57a7330](https://github.com/soerenschneider/fetcharr/commit/57a73305981c72131efe26e7c93a7bb56cc192aa))
* **deps:** bump github.com/segmentio/kafka-go from 0.4.43 to 0.4.44 ([772de24](https://github.com/soerenschneider/fetcharr/commit/772de24bc333d818eafbe7d801f22dcdabf71796))
* **deps:** bump golang from 1.21.2 to 1.21.3 ([3377060](https://github.com/soerenschneider/fetcharr/commit/337706014cd140be8c40f12579ce63c20885b086))

## [1.3.0](https://github.com/soerenschneider/fetcharr/compare/v1.2.0...v1.3.0) (2023-10-13)


### Features

* add nicely formatted filesizes to struct ([73bedee](https://github.com/soerenschneider/fetcharr/commit/73bedee95d6d6d6d0c4e4e52eb4020c3ce34c40f))
* use templates to use in webhook's payload ([6df82cd](https://github.com/soerenschneider/fetcharr/commit/6df82cd12a91f905d36cfe8f98218c901c4082bc))


### Bug Fixes

* check err before trying to close body ([20066ff](https://github.com/soerenschneider/fetcharr/commit/20066ffaa2ecc0642312175cc86a075e97d1ec7b))
* fix duplicate yaml tag ([5bec101](https://github.com/soerenschneider/fetcharr/commit/5bec101fe519bb51953886e55c3ced072a541f34))
* fix printing version information ([5c9bd15](https://github.com/soerenschneider/fetcharr/commit/5c9bd15a718996eeaefc13482f606c5693ba22b4))
* return encoded data even if there's no (valid) stats object ([6edf7a3](https://github.com/soerenschneider/fetcharr/commit/6edf7a30effc1a6715a449eb7638298a0dabd571))

## [1.2.0](https://github.com/soerenschneider/fetcharr/compare/v1.1.1...v1.2.0) (2023-10-13)


### Features

* add support for webhooks ([3cfb4ef](https://github.com/soerenschneider/fetcharr/commit/3cfb4ef2f3ba3bf9677971a7f4b7f4cc0d322587))


### Bug Fixes

* add message to fatal statement ([dce6974](https://github.com/soerenschneider/fetcharr/commit/dce6974f272b9d09381dab87b98a863f89712138))
* **deps:** bump golang from 1.21.1 to 1.21.2 ([baa20c0](https://github.com/soerenschneider/fetcharr/commit/baa20c0d114b0568cc49b7537fe6e8e8ae1f5532))
* **deps:** bump golang.org/x/net from 0.10.0 to 0.17.0 ([5a44dcd](https://github.com/soerenschneider/fetcharr/commit/5a44dcd924fe272137799e5532a9d9b5dadde01b))
* use correct subsystem prefix ([d2b191f](https://github.com/soerenschneider/fetcharr/commit/d2b191f5a67f68e6926efd92ab58a3487c81ebfc))

## [1.1.1](https://github.com/soerenschneider/fetcharr/compare/v1.1.0...v1.1.1) (2023-10-04)


### Bug Fixes

* collect correct runtime ([6239f27](https://github.com/soerenschneider/fetcharr/commit/6239f27b7a229ed5da8d5869665c2ed839c9a581))
* **deps:** bump alpine from 3.18.3 to 3.18.4 ([f7914a2](https://github.com/soerenschneider/fetcharr/commit/f7914a246cbe73f4d8a40e8216338c5967557e42))
* **deps:** bump github.com/go-playground/validator/v10 ([6ad1c3d](https://github.com/soerenschneider/fetcharr/commit/6ad1c3d902b591b47e000fff20002361d5edd673))
* **deps:** bump github.com/prometheus/client_golang ([e73bf81](https://github.com/soerenschneider/fetcharr/commit/e73bf81aed3d71ab0703d1e5f93dc3b6b1439f07))
* **deps:** bump github.com/segmentio/kafka-go from 0.4.42 to 0.4.43 ([8f69439](https://github.com/soerenschneider/fetcharr/commit/8f694393d07129045e60060d18001d1e2095a053))
* rename flag to print version ([9869f0b](https://github.com/soerenschneider/fetcharr/commit/9869f0b534386b54d32821c690cb3450aeca09c2))

## [1.1.0](https://github.com/soerenschneider/fetcharr/compare/v1.0.0...v1.1.0) (2023-09-25)


### Features

* add config option for metrics handler ([727ac2c](https://github.com/soerenschneider/fetcharr/commit/727ac2cf45f214a6474538a2703b57a36da19a09))
* add further metrics ([1746e3f](https://github.com/soerenschneider/fetcharr/commit/1746e3fb9def28b07899ae2daef8b7eabc11fc1f))
* kafka client certs support ([905b7bb](https://github.com/soerenschneider/fetcharr/commit/905b7bb3d9b0a69917e2b31a0ffa68a97c8659a0))


### Bug Fixes

* **deps:** bump github.com/rs/zerolog from 1.30.0 to 1.31.0 ([8defd76](https://github.com/soerenschneider/fetcharr/commit/8defd7692e221d09eccfbc1c6f7fce1101ed3aa0))
* fix parsing stats for linux ([ae9a8ae](https://github.com/soerenschneider/fetcharr/commit/ae9a8ae54ba42b524c49af7ecb5314f5954c7a2c))

## 1.0.0 (2023-09-22)


### Miscellaneous Chores

* release 1.0.0 ([8e55f50](https://github.com/soerenschneider/fetcharr/commit/8e55f50064fb6e19743e66a9c28cb690301838de))
