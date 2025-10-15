import uuid

from aiogram.types import LabeledPrice

from utils.misc.pay_load import Pay


P180 = Pay(
    title="Поповнення на 180 грн",
    description="Поповнення рахунку на 180 гривень",
    currency="UAH",
    prices=[
        LabeledPrice(
            label="180",
            amount=180_00
        )
    ],
    start_parameter=f"create_invoice_180_{uuid.uuid4()}"
)

P1080 = Pay(
    title="Поповнення на 1080 грн",
    description="Поповнення рахунку на 1080 гривень",
    currency="UAH",
    prices=[
        LabeledPrice(
            label="1080",
            amount=1080_00
        )
    ],
    start_parameter=f"create_invoice_1080_{uuid.uuid4()}"
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
    start_parameter=f"create_invoice_200_{uuid.uuid4()}"
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
    start_parameter=f"create_invoice_1200_{uuid.uuid4()}"
)

P350 = Pay(
    title="Поповнення на 350 грн",
    description="Поповнення рахунку на 350 гривень",
    currency="UAH",
    prices=[
        LabeledPrice(
            label="350",
            amount=350_00
        )
    ],
    start_parameter=f"create_invoice_350_{uuid.uuid4()}"
)

P2100 = Pay(
    title="Поповнення на 2100 грн",
    description="Поповнення рахунку на 2100 гривень",
    currency="UAH",
    prices=[
        LabeledPrice(
            label="2100",
            amount=2100_00
        )
    ],
    start_parameter=f"create_invoice_2100_{uuid.uuid4()}"
)
