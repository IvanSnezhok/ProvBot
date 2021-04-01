def format_number(telegram_number: str):
    telegram_number = telegram_number.replace('+', '')
    telegram_number = telegram_number.replace('(', '')
    telegram_number = telegram_number.replace(')', '')
    telegram_number = telegram_number[2:]
    r = telegram_number[:3]
    r1 = telegram_number[3:6]
    r2 = telegram_number[6:]
    r2 = r2[:2]
    r3 = telegram_number[8:]
    result = "{}-{}-{}-{}".format(r, r1, r2, r3)
    return result

number = "+380(93)9555270"
print(format_number(number))
