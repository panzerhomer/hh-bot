import logging
import requests
import psycopg2

from config import _TOKEN
from sql_queries import _CREATE_TABLE_QUERY, _DROP_TABLE_QUERY, _DELETE_DUPLICATES_QUERY


def create_table(connection: psycopg2.connect):
    cursor = connection.cursor()
    cursor.execute(_CREATE_TABLE_QUERY)
    connection.commit()
    cursor.close()
    logging.info("Создана таблица с вакансиями (vacancies)")


def drop_table(connection: psycopg2.connect):
    cursor = connection.cursor()
    cursor.execute(_DROP_TABLE_QUERY)
    connection.commit()
    cursor.close()
    logging.info("Удалена таблица с вакансиями (vacancies)")


def get_hh_vacancies(city_name: str, vacancy: str, page: int) -> dict:
    url = 'https://api.hh.ru/vacancies'
    params = {
        'text': f"{vacancy} {city_name}",
        'area': city_name,
        'specialization': 1,
        'per_page': 100,
        'page': page
    }
    headers = {
        'Authorization': f'Bearer {_TOKEN}'
    }

    response = requests.get(url, params=params, headers=headers)
    response.raise_for_status()
    return response.json()


def get_skills(job_id: str) -> str:
    url = f'https://api.hh.ru/vacancies/{job_id}'
    headers = {
        'Authorization': f'Bearer {_TOKEN}'
    }
    response = requests.get(url, headers=headers)
    response.raise_for_status()
    return ', '.join([skills['name'] for skills in response.json().get('key_skills', [])])


def get_industry_name(company_id):
    if company_id is None:
        return 'Unknown'

    url = f'https://api.hh.ru/employers/{company_id}'
    response = requests.get(url)
    if response.status_code == 404:
        return 'Unknown'
    response.raise_for_status()

    data = response.json()
    if 'industries' in data and len(data['industries']) > 0:
        return data['industries'][0].get('name')

    return 'Unknown'


def delete_duplicates(connection: psycopg2.connect):
    cursor = connection.cursor()
    cursor.execute(_DELETE_DUPLICATES_QUERY)
    connection.commit()
    cursor.close()

    logging.info("Дубликаты в таблице ваканский ('vacancies') удалены")
