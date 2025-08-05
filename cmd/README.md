# UI Bakery

This folder contains CLI commands, subcommands, flags and similar.

## Commands

### Queue

The command queue performs operations over the message queue.

- claim the records

  ```bash
  uibakery queue claim [-size 10]
  ```

  Claims up to `size` free records from the database queue. Records must be in
  appropriate state (`NEW` or `PENDING` or `GRIEF` - for the last one
  additionally records must rest for at least 15 minutes) and unclaimed prior
  the operation.

- unclaim the claimed records

  ```bash
  uibakery queue unclaim --claim-id 12345678-9abc-def0-1234-56789abcdef0
  ```

  Flag `claim-id` is mandatory. Command will display number of records
  processed.

- count records

  ```bash
  $ uibakery queue unclaim --count
              | cryplex             | shopify             | shopify_sync        | zendesk             | total               |
              | claimed  : free     | claimed  : free     | claimed  : free     | claimed  : free     | claimed  : free     |
  ---------------------------------------------------------------------------------------------------------------------------
  CANCELED-GF |        0 :       23 |          :          |          :          |          :          |        0 :       23 |
  CLOSED-GF   |        0 :       18 |          :          |          :          |          :          |        0 :       18 |
  CLOSED-NA   |        0 :      604 |          :          |          :          |        0 :      936 |        0 :     1540 |
  GRIEF       |        0 :        7 |          :          |          :          |          :          |        0 :        7 |
  NEW         |       19 :        0 |        0 :     1038 |          :          |          :          |       19 :     1038 |
  PROCESSED   |          :          |          :          |        0 :     9331 |          :          |        0 :     9331 |
  SUCCESS     |        0 :     1095 |          :          |          :          |          :          |        0 :     1095 |
  SUCCESS-CR  |          :          |          :          |          :          |        0 :        7 |        0 :        7 |
  SUCCESS-UD  |          :          |          :          |          :          |        0 :      363 |        0 :      363 |
  ```

- displays claims

  ```bash
  uibakery queue list-claims [--status NEW]
  ```

  Command lists claims in the database, along with the number of records claimed.
  User can filter to only a specified submission status (`NEW`, for example) and
  this will also affect the number of record if the claimed records are in
  different state.

## Flags

- global
  - hostname
  - port
  - username
  - secret
  - database
- queue claim
  - size
- queue unclaim
  - claim-id
- queue list-claims
  - status
