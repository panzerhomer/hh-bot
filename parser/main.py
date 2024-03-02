import logging
import random
import time

import pandas as pd
import psycopg2
import requests

from sqlalchemy import create_engine, Engine
from typing import Dict, Any

from config import _DATABASE, _USER, _PASSWORD, _HOST, _PORT
from parser_settings import _CITIES, _SPECIALIZATIONS, _MAX_PAGE
from functions import create_table, drop_table, get_hh_vacancies, get_skills, get_industry_name, delete_duplicates

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')


def parse_page(engine: Engine,
               hh_data: Dict[str, Any],
               specialization: str,
               city_name: str):
    for vacancy in hh_data['items']:
        if specialization.lower() not in vacancy['name'].lower():
            continue

        info = {
            'city': [city_name],
            'company': [vacancy['employer']['name']],
            'industry': [get_industry_name(vacancy['employer'].get('id'))],
            'title': [f"{vacancy['name']} ({city_name})"],
            'keywords': [vacancy['snippet'].get('requirement', '')],
            'skills': [get_skills(vacancy['id'])],
            'experience': [vacancy['experience'].get('name', '')],
            'salary': vacancy['salary'],
            'url': [vacancy['alternate_url']]
        }
        info['salary'] = "з/п не указана" if info['salary'] is None else info['salary'].get('from', '')

        df = pd.DataFrame(info)
        df.to_sql('vacancies', engine, if_exists='append', index=False)


def parse_hh_vacancies(connection: psycopg2.connect):
    engine = create_engine(
        f'postgresql://{_USER}:{_PASSWORD}@{_HOST}:{_PORT}/{_DATABASE}')

    drop_table(connection)
    create_table(connection)

    for city_name, city_id in _CITIES.items():
        for specialization in _SPECIALIZATIONS:
            page = 0
            hh_data = get_hh_vacancies(city_id, specialization, page)
            while hh_data.get('items') and page < _MAX_PAGE:
                try:
                    parse_page(engine, hh_data, specialization, city_name)
                except requests.HTTPError as e:
                    logging.error(f"Возникла ошибка в городе:{city_name}: {e}")
                else:
                    page += 1
                    time.sleep(random.uniform(3, 6))
                    hh_data = get_hh_vacancies(city_id, specialization, page)

        connection.commit()

    logging.info("Окончание прасинга")


def start_parse_hh_vacancies():
    logging.info("Запуск парсера")
    try:
        conn_string = f"host={_HOST} dbname={_DATABASE} user={_USER} password={_PASSWORD}"
        with psycopg2.connect(conn_string) as connection:
            parse_hh_vacancies(connection)
            delete_duplicates(connection)
    except Exception as e:
        logging.error(f"Ошибка при парсинге: {e}")


start_parse_hh_vacancies()