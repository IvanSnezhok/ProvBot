def format_number(telegram_number: str):
    telegram_number = "".join(num for num in telegram_number if num not in " +()")[2:]
    r = telegram_number[:3]
    r1 = telegram_number[3:6]
    r2 = telegram_number[6:]
    r2 = r2[:2]
    r3 = telegram_number[8:]
    result = "{}-{}-{}-{}".format(r, r1, r2, r3)
    return result
