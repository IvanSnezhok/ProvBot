from aiogram.types import LabeledPrice

from utils.misc.pay_load import Pay

P150 = Pay(
    title="Поповнення на 150 грн",
    description="Поповнення рахунку на 150 гривень",
    currency="UAH",
    prices=[
        LabeledPrice(
            label="150",
            amount=150_00
        )
    ],
    start_parameter="create_invoice_150"
)

P900 = Pay(
    title="Поповнення на 900 грн",
    description="Поповнення рахунку на 900 гривень",
    currency="UAH",
    prices=[
        LabeledPrice(
            label="900",
            amount=900_00
        )
    ],
    start_parameter="create_invoice_900"
)

P200 = Pay(
    title="Поповнення на 200 грн",
    description="Поповнення рахунку на 200 гривень",
    currency="UAH",
    prices=[
        LabeledPrice(
            label="200",
            amount=200_00
        )
    ],
    start_parameter="create_invoice_200"
)

P1200 = Pay(
    title="Поповнення на 1200 грн",
    description="Поповнення рахунку на 1200 гривень",
    currency="UAH",
    prices=[
        LabeledPrice(
            label="1200",
            amount=1200_00
        )
    ],
    start_parameter="create_invoice_1200"
)
