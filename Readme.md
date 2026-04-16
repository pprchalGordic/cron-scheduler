# Návod na použití Plánovače úloh

## Co to je
Tato aplikace umožňuje automaticky spouštět různé úlohy (skripty) ve vámi zadaný čas a dny. Je to obdoba cron služby z Linuxu, ale pro Windows.

## Jak to funguje
- Aplikace se spouští v pravidelných intervalech (doporučuje se každých 30 minut)
- Zkontroluje, zda má nějakou úlohu spustit v aktuální čas
- Spustí úlohu pouze jednou za den (i když se aplikace spustí vícekrát)
- Uloží výstup každé úlohy do logu


## Nastavení konfigurace
Vytvořte nebo upravte soubor `schedule.yaml` (nebo `config.yaml`) podle vzoru:

```yaml
jobs:
  - name: Rotace logů  
    run_at: "8:00"
    command: logrotate C:\Log C:\LogArchive

  - name: Zálohování
    run_at: "23:00"
    days: [Mon, Tue, Wed, Thu, Fri]
    command: "-NoProfile -File C:\Scripts\backup.ps1"
```

## Nastavení úloh

### Parametry úlohy:
- **name**: Název úlohy (pro rozpoznání v lozích)
- **run_at**: Čas spuštění ve formátu "HH:MM" (např. "08:30", "23:15")
- **days**: Dny spuštění (nepovinné)
  - Možné hodnoty: Mon, Tue, Wed, Thu, Fri, Sat, Sun
  - Pokud není uvedeno, úloha se spustí každý den
- **command**: Příkaz ke spuštění

### Druhy příkazů:

#### PowerShell skripty:
```yaml
command: "-NoProfile -File C:\Scripts\muj_skript.ps1"
```

#### Rotace logů (speciální funkce):
```yaml
command: logrotate C:\SourceFolder C:\ArchiveFolder
```
- Automaticky zazipuje a přesune staré log složky
- Zpracovává pouze složky starší než včerejšek
- Složky musí mít formát názvu: YYYY_MM_DD nebo YYYY-MM-DD

## Nastavení ve Windows

### Automatické spouštění přes Plánovač úloh Windows:

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
    - Zaškrtněte "Repeat task every:" 30 minutes
    - For a duration of: 1 day

## Složky a soubory

Po spuštění aplikace se vytvoří tyto složky a soubory:

- **logs/**: Obsahuje logy z každé úlohy
  - Každá úloha má svůj soubor: `název_úlohy.log`
- **locks/**: Dočasné soubory zamykání (automaticky se mažou)
- **state.json**: Sleduje, které úlohy už byly dnes spuštěny

## Příklady použití

### Denní zálohování dat:
```yaml
- name: Backup dat
  run_at: "02:00"
  command: "-NoProfile -File C:\Scripts\backup_data.ps1"
```

### Úklid tempů pouze v pracovní dny:
```yaml
- name: Cleanup temp
  run_at: "18:00"
  days: [Mon, Tue, Wed, Thu, Fri]
  command: "-NoProfile -File C:\Scripts\cleanup_temp.ps1"
```

### Rotace logů aplikace:
```yaml
- name: Rotace aplikačních logů
  run_at: "01:00"
  command: logrotate C:\AppLogs C:\LogArchive
```

## Řešení problémů

### Aplikace se nespustí:
- Zkontrolujte, že máte nainstalované Go
- Zkompilujte znovu: `go build`

### Úloha se nespustila:
1. Zkontrolujte čas a formát v `run_at`
2. Ověřte, že dny jsou správně zadané
3. Podívejte se do logu úlohy v složce `logs/`

### PowerShell skripty se nespouští:
- Zkontrolujte cesty k souborům
- Ujistěte se, že máte oprávnění ke spuštění skriptů
- Otestujte skript ručně v PowerShellu

### Rotace logů nefunguje:
- Ověřte, že složky existují
- Zkontrolujte, že source složky mají správný formát názvu (YYYY_MM_DD)
- Ujistěte se, že máte práva k zápisu do cílové složky

## Zobrazení výstupu

Po spuštění aplikace bez parametrů se zobrazí nápověda:
```
process-cron launcher
schedule in task manager, configure by schedule.yaml

runs every n minutes (30) recommended and executes scripts
```

Všechny výstupy úloh se ukládají do souborů v složce `logs/`.