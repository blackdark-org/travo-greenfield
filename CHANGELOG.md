# Changelog

## 0.1.0 (2026-11-14)


### Features

* add Device port with mock and CLI adapters, guards, atomic files, lifecycle ([8a0982d](https://github.com/blackdark-org/travo-greenfield/commit/8a0982d905f291320b67c96e0b2176037bcd1107))
* add Device port with mock and CLI adapters, guards, atomic files, lifecycle ([58565c1](https://github.com/blackdark-org/travo-greenfield/commit/58565c16120570d8778642d0d4f562331d1d8d77))
* add DNS layer stack with fixed precedence and guarded apply ([#6](https://github.com/blackdark-org/travo-greenfield/issues/6)) ([a48383b](https://github.com/blackdark-org/travo-greenfield/commit/a48383b50ded78dfa6ac6efe624eaa989bee657e))
* add PackageManager port with apk/opkg adapters, --mock wiring and size gate ([912990f](https://github.com/blackdark-org/travo-greenfield/commit/912990f802a67128f1fb022dfa919c73daf46ce8))
* **api:** add auth, first-run and transaction endpoints with default-deny ([09f4d01](https://github.com/blackdark-org/travo-greenfield/commit/09f4d0198da7638a4ffc1e05caa3882f9247dcd2))
* **auth:** add session store, epoch file, rate limiter and password rules ([c9ad121](https://github.com/blackdark-org/travo-greenfield/commit/c9ad121722c2d30a30fb7beeec9fd440a7fab123))
* **daemon:** serve the web UI with gzip, caching and CSP ([2a380bf](https://github.com/blackdark-org/travo-greenfield/commit/2a380bf934aba1e6677b4e09f0dc6a363bbe27aa)), closes [#18](https://github.com/blackdark-org/travo-greenfield/issues/18)
* **dashboard:** add shared status event stream ([7c9e443](https://github.com/blackdark-org/travo-greenfield/commit/7c9e443366939ce645ebdbff04a7e5784424a060))
* **dashboard:** share live status subscription ([1e5f364](https://github.com/blackdark-org/travo-greenfield/commit/1e5f364103700a80049a5714ab5abdb59c374a77))
* **dashboard:** show router resource summary ([5b51f17](https://github.com/blackdark-org/travo-greenfield/commit/5b51f174362960a3652260e1ce8b2baca8fa08d9))
* **dashboard:** show throughput by interface ([f9c856e](https://github.com/blackdark-org/travo-greenfield/commit/f9c856e4b2d835643bfe2adaa8586c4247726edb))
* **dashboard:** show uplink and internet status ([12a8041](https://github.com/blackdark-org/travo-greenfield/commit/12a8041cf63b64cc11e8ced4727dd46b50b79375))
* **dev:** add mise tasks for the mock dev server and screenshots ([1531c3c](https://github.com/blackdark-org/travo-greenfield/commit/1531c3c136dab5a2b15cfd01565f0c7a6b073f1d)), closes [#17](https://github.com/blackdark-org/travo-greenfield/issues/17)
* **device:** add Reload to the Device port and read uci via export ([57e2d91](https://github.com/blackdark-org/travo-greenfield/commit/57e2d912c3a5b71e6eb6385504925dc64183580e))
* **device:** add root account port (verify, detect, set password) ([574b92a](https://github.com/blackdark-org/travo-greenfield/commit/574b92a6280b1223132762534f005560ce33be49))
* **device:** parse netifd, neighbour and hostapd payloads ([dff3679](https://github.com/blackdark-org/travo-greenfield/commit/dff3679361ed914879010de88c112f1599b609e2))
* **device:** read WireGuard handshake evidence ([676793e](https://github.com/blackdark-org/travo-greenfield/commit/676793e494f82ac022c2228d21bc870eefec89d4))
* **lockout:** classify callers and refuse lockout-risk changes ([2c65911](https://github.com/blackdark-org/travo-greenfield/commit/2c659117183d3e9d826413c65278c77873692583))
* **openapi:** add command that prints the generated API spec ([451541a](https://github.com/blackdark-org/travo-greenfield/commit/451541ab0bb8a03ce3d029b3a1789285dfa56424))
* **packaging:** add the release tarball, init script, install and uninstall ([9f83cdf](https://github.com/blackdark-org/travo-greenfield/commit/9f83cdf3cd0015c7b0e446fc3f222016d3f94c55)), closes [#19](https://github.com/blackdark-org/travo-greenfield/issues/19)
* **release:** sign the release tarball and verify it on install ([22a6492](https://github.com/blackdark-org/travo-greenfield/commit/22a6492da79463f7b7457b79b36cb538c07bc154)), closes [#21](https://github.com/blackdark-org/travo-greenfield/issues/21)
* **runner:** run commands with stdin and add periodic lifecycle workers ([2dfdf2d](https://github.com/blackdark-org/travo-greenfield/commit/2dfdf2d275a6e229cb6d78fd0f01c2fb19937c79))
* **system:** apply signed updates with guarded rollback ([5b0b7d3](https://github.com/blackdark-org/travo-greenfield/commit/5b0b7d37651c1a31f95e4327167dccffceaaec24))
* **system:** check for stable releases on demand ([b1203fd](https://github.com/blackdark-org/travo-greenfield/commit/b1203fdec71e8c4e28a83fea6cc38f920ad7005d))
* **system:** download router configuration backups ([994912b](https://github.com/blackdark-org/travo-greenfield/commit/994912b945b439398550cd5b8d9bde9d90bfda5a))
* **system:** show device info on the System page ([43e50a2](https://github.com/blackdark-org/travo-greenfield/commit/43e50a24d2cd68e1c552dfe1fcdcffed569d3b7d)), closes [#50](https://github.com/blackdark-org/travo-greenfield/issues/50)
* **system:** verify and stage signed release packages ([e08bbcb](https://github.com/blackdark-org/travo-greenfield/commit/e08bbcb2969dce92f7a01bb95e3cf01ac9fd82a0))
* **trust:** allow admin only from LAN subnets ([a989dd1](https://github.com/blackdark-org/travo-greenfield/commit/a989dd1c0725a4cadc9db41ff3cecc907f991f3f))
* **txn:** add transactions with detached watchdog and startup recovery ([01e4f92](https://github.com/blackdark-org/travo-greenfield/commit/01e4f92c186ad14c16a579f275e4e1c1909fcedd))
* **uciwrite:** add per-config UCI locks ([515b020](https://github.com/blackdark-org/travo-greenfield/commit/515b0205487b4a33ccad5371bc2c1986d820fe16))
* **vpn:** expose WireGuard tunnel status ([9e47184](https://github.com/blackdark-org/travo-greenfield/commit/9e47184903177ae21ab8120761c1eb76c945611b))
* **vpn:** manage WireGuard profiles ([d415491](https://github.com/blackdark-org/travo-greenfield/commit/d4154911ccb36aecebdc5b130a7d5fc2ed86c490))
* **web:** add login, first-run setup and session handling ([30c1c8e](https://github.com/blackdark-org/travo-greenfield/commit/30c1c8e0022353c000b580e2e21855f4ae0f2f9c)), closes [#16](https://github.com/blackdark-org/travo-greenfield/issues/16)
* **web:** add Svelte SPA skeleton with generated API client ([dba66e2](https://github.com/blackdark-org/travo-greenfield/commit/dba66e224ef6cb8e0f2e1854b3155eba8e00b873))
* **web:** show WiFi summary and scan results ([5249e3c](https://github.com/blackdark-org/travo-greenfield/commit/5249e3c512c37c6f60a38db6b03f8f886d9e59bf))
* **web:** show WireGuard tunnel status ([c1436ce](https://github.com/blackdark-org/travo-greenfield/commit/c1436ce9a6b05accaab1b853e960e0a71c63f477))
* **wifi:** add read-only summary and scan API ([cda78fb](https://github.com/blackdark-org/travo-greenfield/commit/cda78fb87c2303bdef819877763a36e99f671faf))


### Bug Fixes

* **api:** drain dashboard producers before restart ([bf36dc1](https://github.com/blackdark-org/travo-greenfield/commit/bf36dc10f29b879eaf77023640b5c992c8d02ed7))
* **api:** refuse unexpected Host headers (DNS rebinding) ([75d9085](https://github.com/blackdark-org/travo-greenfield/commit/75d9085d835dc47e067c89f78c836428686f59ae))
* **app:** remove the mock state dir on shutdown and drop unused session field ([0fd4c03](https://github.com/blackdark-org/travo-greenfield/commit/0fd4c03be6f14983a29fbb86326f4fc155b9a432))
* **dashboard:** resume events after bfcache restore ([7ed9009](https://github.com/blackdark-org/travo-greenfield/commit/7ed90091ce6763f03caf705d24380915059098aa))
* **device:** distinguish transient uci export failures ([efe7a6d](https://github.com/blackdark-org/travo-greenfield/commit/efe7a6d22f014154e1cac1897b835593a15b6f0e))
* **device:** match complete command diagnostics ([d677972](https://github.com/blackdark-org/travo-greenfield/commit/d677972e7a2ba4feef3c6c7242e8c212680751b5))
* **device:** reject reserved WireGuard selectors ([0a1eed5](https://github.com/blackdark-org/travo-greenfield/commit/0a1eed53ec91a456dc6d0b290757b86e1b3e4013))
* **openapi:** build the spec with the new api.Deps ([07b20ad](https://github.com/blackdark-org/travo-greenfield/commit/07b20addfd168d9bec769e5b38b7a4cc2f8d07dc))
* **packaging:** lowercase the hostname with explicit ranges for BusyBox tr ([3892ada](https://github.com/blackdark-org/travo-greenfield/commit/3892adaa7f8073bc67e92351c1e1cd44477c9af6)), closes [#19](https://github.com/blackdark-org/travo-greenfield/issues/19)
* **packaging:** move LuCI for reversible port 80 installs ([8ebeefd](https://github.com/blackdark-org/travo-greenfield/commit/8ebeefd8b1f780bd1b7d79d415e25910ea076b2c))
* **packaging:** target fresh OpenWrt installations ([a8f8f33](https://github.com/blackdark-org/travo-greenfield/commit/a8f8f338d8a60da6ed012d2571de72bcfdcfc02f))
* **runner:** stop command descendants on cancellation ([5264217](https://github.com/blackdark-org/travo-greenfield/commit/5264217cecab2a825dc6eff482aae957597e08e6))
* **txn:** enforce shared live apply deadline ([d7cf4e1](https://github.com/blackdark-org/travo-greenfield/commit/d7cf4e1a25eb8ab43b02766351432b50a03c9d6f))
* **vpn:** expire stale status on monotonic time ([ae4a370](https://github.com/blackdark-org/travo-greenfield/commit/ae4a370319b998d1756017cdac01e63ca907e3ec))
* **web:** disambiguate duplicate scan access points ([12f7076](https://github.com/blackdark-org/travo-greenfield/commit/12f70760591addb11e5c85432d44b50c048d1a48))
* **web:** keep polling refreshes from announcing loading ([00eecde](https://github.com/blackdark-org/travo-greenfield/commit/00eecded8a01bb09b88507b32b168c8d66057690))
* **web:** recover polling from slow and hung requests ([5b532f6](https://github.com/blackdark-org/travo-greenfield/commit/5b532f6e435a59fb7389fb979a157c8e47676038))
* **wifi:** cancel scans with daemon lifecycle ([7eb9d44](https://github.com/blackdark-org/travo-greenfield/commit/7eb9d44b1e4ebae4524b0efee6297e188043b637))


### Refactoring

* **api:** derive VPN status schema name from Go type ([7f410ba](https://github.com/blackdark-org/travo-greenfield/commit/7f410ba6e5b52ea2cf31cee699ae64823ef11345))
* **dashboard:** extract app event ownership ([5928436](https://github.com/blackdark-org/travo-greenfield/commit/5928436b3ae4f54bfc49e6f7bef2a39456c8eee6))
* **dnslayer:** apply through txn.Manager.Apply ([46a26b1](https://github.com/blackdark-org/travo-greenfield/commit/46a26b13c8b0218dc3dcafa219bd5ba84d4e59ae))
* **packaging:** share OpenWrt container test helpers ([7cef0ff](https://github.com/blackdark-org/travo-greenfield/commit/7cef0ffeed3cd04afec880884a82cca6551385bc))


### Documentation

* add 0.1.0 spec, invariants, ADRs and agent setup ([6e0cb7f](https://github.com/blackdark-org/travo-greenfield/commit/6e0cb7fc0786f5a6d9813cde46914cd6fd9963f8))
* add handoff notes for the next session ([4d0daab](https://github.com/blackdark-org/travo-greenfield/commit/4d0daab1ecc4d81eb0fed4ef44f4dfdf3d34ac29))
* **adr:** document split UI size budgets ([76f56b3](https://github.com/blackdark-org/travo-greenfield/commit/76f56b32cc8329e0473022a396005441c4423c13))
* move planning to GitHub issues and the project board ([083796f](https://github.com/blackdark-org/travo-greenfield/commit/083796f31e09ae044d13520cc5d506378775ffa3))
* move the glossary to GLOSSARY.md at the root ([5bdc5fb](https://github.com/blackdark-org/travo-greenfield/commit/5bdc5fb08bc2d570422bd43387683179c3cb8e22))
* note why old Travo can run TypeScript 7 and we cannot ([95b2dcf](https://github.com/blackdark-org/travo-greenfield/commit/95b2dcf6dc908ed625febf6161eb32f8ebfbeaaa))
* record auth implementation and the loopback trust rule ([53424e4](https://github.com/blackdark-org/travo-greenfield/commit/53424e4ef4411381751b1f7052cfc7bd60c48b18))
* record device facts, testing layers and spike outcomes ([f51d181](https://github.com/blackdark-org/travo-greenfield/commit/f51d181a7a59b0a9be4dc05ef5c13e368c0d6af9))
* record web stack details, gen workflow and UI checks ([b28bc39](https://github.com/blackdark-org/travo-greenfield/commit/b28bc392c2770a9bc21914572b343789798484f9))
* retire the slice-3 handoff and refresh handoff status ([8fe5cdb](https://github.com/blackdark-org/travo-greenfield/commit/8fe5cdbe916c9a88ac514b70cd5899431ed4268b))
* set up agent skills config and the issue-to-merge workflow ([fd94645](https://github.com/blackdark-org/travo-greenfield/commit/fd94645319283332b4c7639a9b89a29fbae34cf3))
* use the travo-greenfield state and guard paths ([30c05cb](https://github.com/blackdark-org/travo-greenfield/commit/30c05cb2061e7ceec3c44b36cbafc06f80cbd27b))
* **vpn:** define tunnel status semantics ([bc97d5c](https://github.com/blackdark-org/travo-greenfield/commit/bc97d5c763ca699c9270a4072cc55e910f492811))
* **vpn:** name tunnel status endpoint ([6795816](https://github.com/blackdark-org/travo-greenfield/commit/679581621c509b597cc0bbcbe3e9396b2dd6c10e))
