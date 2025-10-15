import xml.etree.ElementTree as ET
import requests


def get_products():
    url = "https://sunrise.co.ua/marketplace-integration/google-feed/27bd50b19dc9c0d7d1ed3c5a7278cc49?langId=3"
    response = requests.get(url)
    root = ET.fromstring(response.content)

    products = []
    for item in root.findall('.//item'):
        try:
            product = {
                'id': item.find('g:id', namespaces={'g': 'http://base.google.com/ns/1.0'}).text,
                'title': item.find('g:title', namespaces={'g': 'http://base.google.com/ns/1.0'}).text,
                'description': item.find('g:description', namespaces={'g': 'http://base.google.com/ns/1.0'}).text,
                'link': item.find('g:link', namespaces={'g': 'http://base.google.com/ns/1.0'}).text,
                'image_link': item.find('g:image_link', namespaces={'g': 'http://base.google.com/ns/1.0'}).text,
                'price': item.find('g:price', namespaces={'g': 'http://base.google.com/ns/1.0'}).text,
                'category': item.find('g:product_type', namespaces={'g': 'http://base.google.com/ns/1.0'}).text,
            }
            
            # Додаємо додаткові зображення, якщо вони є
            additional_images = item.findall('g:additional_image_link', namespaces={'g': 'http://base.google.com/ns/1.0'})
            if additional_images:
                product['additional_images'] = [img.text for img in additional_images]
                
            products.append(product)
        except AttributeError as e:
            print(f"Помилка при парсингу товару: {e}")
            continue

    return products
