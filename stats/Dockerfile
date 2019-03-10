FROM python:3.7.2

ENV PYTHONUNBUFFERED 1
RUN mkdir /code
WORKDIR /code

RUN pip install poetry
ADD poetry.lock .
ADD pyproject.toml .
RUN poetry install
