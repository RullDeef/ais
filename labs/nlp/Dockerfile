FROM python:3.10-alpine

RUN pip install --upgrade pip && \
    pip install nltk pymorphy2

# install additional nlp modules
RUN python -c "import nltk; nltk.download('stopwords'); nltk.download('punkt')"

WORKDIR /app
COPY processor.py /app
CMD ["python", "processor.py"]
