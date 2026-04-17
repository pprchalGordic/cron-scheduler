# Návod na použití plánovače úloh

## Co to je
Pomocná aplikace při implementace be-databázových GINIS modulů (APG/DKS/AIB).
Umožňuje automaticky spouštět různé úlohy (skripty) ve vámi zadaný čas a dny. Je to obdoba cron služby z Linuxu, ale pro Windows.

Umožní přenést konfiguraci více naplánovaných úloh v jediném souboru.
Časová přesnost není pro velké systémy důležitá, postačuje spustit v dobu, kdy je malé zatížení (například v 3:00 ráno).

## Jak to funguje
- `-cron` aplikace se spouští v pravidelných intervalech (doporučuje se každých 30 minut přes TaskManager)
- `-encrypt` zašifruje heslo do `dpapi:tvaru`
- `-decrypt` rozšifruje heslo z `dpapi:tvaru`
- `-logrotate` provede archivaci složek logů

## -cron
1) zkontroluje, zda má nějakou úlohu spustit v aktuální čas
2) spustí úlohu pouze jednou za den (i když se aplikace spustí vícekrát)
3) uloží výstup každé úlohy do logu

Vytvořte nebo upravte soubor `config.yaml` podle vzoru:

```yaml
jobs:
  - name: Rotace logů  
    run_at: "8:00"
    command: logrotate C:\Log C:\LogArchive

  - name: Zálohování
    run_at: "23:00"
    days: [1, 2, 3, 4, 5]
    command: "-NoProfile -File C:\Scripts\backup.ps1"
```

### Popis parametrů:
- **name**: Název úlohy (pro rozpoznání v lozích)
- **run_at**: Čas spuštění ve formátu "HH:MM" (např. "08:30", "23:15")
- **days**: Dny spuštění (nepovinné)
  - Čísla 1–7, kde 1=Po, 2=Út, 3=St, 4=Čt, 5=Pá, 6=So, 7=Ne
  - Pokud není uvedeno, úloha se spustí každý den
- **command**: Příkaz ke spuštění


### Druhy příkazů `command`:

#### PowerShell skripty:
```yaml
command: "-NoProfile -File C:\Scripts\muj_skript.ps1"
```

#### Rotace logů
Automaticky zazipuje a přesune staré log složky.
Zpracovává pouze složky starší než včerejšek.

Může být 
spuštěno i interaktivně (viz `-logrotate`)
> Složky musí mít formát názvu: YYYY_MM_DD nebo YYYY-MM-DD

```yaml
command: logrotate C:\SourceFolder C:\ArchiveFolder
```

## Automatické spouštění přes Plánovač úloh Windows:

1. Otevřete **Task Scheduler** (Plánovač úloh)
2. Klikněte na **"Create Basic Task"** (Vytvořit základní úlohu)
3. Zadejte název: "Cron Scheduler"
4. Vyberte **"Daily"** (Denně)
5. Nastavte čas startu: např. 00:00
6. Vyberte **"Start a program"** (Spustit program)
7. Program: cesta k `cron-scheduler.exe`
8. Argumenty: nechte prázdné
9. Složka: cesta ke složce s aplikací
10. V **pokročilých nastaveních**:
    - Zaškrtněte "Repeat task every:" 15 minutes
    - For a duration of: 1 day

## Složky a soubory

Po spuštění aplikace se vytvoří tyto složky a soubory:

- **logs/**: Obsahuje logy z každé úlohy
  - Každá úloha má svůj soubor: `název_úlohy.log`
- **locks/**: Dočasné soubory zamykání (automaticky se mažou)
- **state.json**: Sleduje, které úlohy už byly dnes spuštěny


## -encrypt
Zašifruje heslo do `dpapi:tvaru` pro GINIS-vault. 

## -decrypt
Rozšifruje heslo z `dpapi:tvaru`.

> buďte opatrní



## Příklady
### spustit naplánované úkoly
```
.\cron-scheduler.exe -cron

```

### zašifrovat heslo do trezoru
```
.\cron-scheduler.exe -encrypt
```

### rozšifrovat heslo z trezoru
interaktivní způsob - zeptá se uživatele na `dpapi:tvar`
> buďte opatrní
```
.\cron-scheduler.exe -decrypt
```

### rozšifrovat heslo z trezoru na windows 2022
v `cesta_k_souboru.txt` musí být zašifrované heslo v `dpapi:tvaru`

> buďte opatrní
```
.\cron-scheduler.exe -decrypt cesta_k_souboru.txt
```
