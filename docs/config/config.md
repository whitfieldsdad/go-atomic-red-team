# Configuration

## go-atomic

| Variable                               | Default Value                        | Description                                                              |
|----------------------------------------|--------------------------------------|--------------------------------------------------------------------------|
| GO_ATOMIC_TEST_DIR                     | data/tests                           | Directory of tests                                                       |
| GO_ATOMIC_TEST_INVOCATION_DIR          | data/test_invocations                | Directory of test invocations                                            |
| GO_ATOMIC_TEST_INVOCATION_STATUS_DIR   | data_test_invocation_statuses        | Directory of test invocation statuses                                    |
| GO_ATOMIC_TEST_INVOCATION_RESULT_DIR   | data/test_invocation_results         | Directory of test invocation results                                     |
| GO_ATOMIC_TEST_INVOCATION_REQUEST_DIR  | data/test_invocation_requests        | Directory of test invocation requests                                    |
| GO_ATOMIC_TEST_INVOCATION_RESPONSE_DIR | data/test_invocation_responses       | Directory of test invocation responses                                   |
| GO_ATOMIC_ENABLE_ARCHIVE_MODE          | false                                | All repositories will be packaged as archives                            |
| GO_ATOMIC_ENABLE_ART                   | true                                 | Commands from Atomic Red Team will be included                           |
| GO_ATOMIC_ENABLE_LOLBAS                | false                                | Commands from LOLBAS will be included                                    |
| GO_ATOMIC_ENABLE_LOLDRIVERS            | false                                | Commands from LOLDRIVERS will be included                                |
| GO_ATOMIC_ENABLE_GTFOBINS              | false                                | Commands from GTFOBINS will be included                                  |
| GO_ATOMIC_ARCHIVE_KEY                  | ae175d6d-d952-4fc6-b967-fc6fa6f61fce | If provided, archives will be encrypted using the provided symmetric key |
| GO_ATOMIC_ART_DIR                      | data/atomic-red-team                 | Path to Atomic Red Team repository                                       |
| GO_ATOMIC_LOLBAS_DIR                   | data/LOLBAS                          | Path to LOLBAS repository                                                |
| GO_ATOMIC_GTFOBINS_DIR                 | data/GTFOBins                        | Path to GTFOBINS repository                                              |
| GO_ATOMIC_LOLDRIVERS_DIR               | data/LOLDRIVERS                      | Path to LOLDRIVERS repository                                            |
