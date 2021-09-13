# ynab-helper

Telegram Bot assistant for YNAB app. 

Parses SMS text like:

```text
Priorbank. Karta 4***1122 10-09-2021 15:40:19. Oplata 12.96 BYN. BLR SHOP COOL NANE.   Spravka: 80123456789
```

and saves it into DB. 

Invalid SMS that cannot be parsed also saved into DB for invesigation purposes.


## Verify test coverage

```bash
go tool cover -html=cover.out
```

opens default web browser with detailed code coverage.
