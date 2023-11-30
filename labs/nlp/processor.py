import os
import pymorphy2
import json
from nltk.corpus import stopwords
from nltk.util import ngrams
from nltk.tokenize import word_tokenize
import http
from http import server
from urllib.parse import urlparse, parse_qs


morph = pymorphy2.MorphAnalyzer()
keep_stopwords = tuple() # ('не', 'всех')
stopwords = list(filter(lambda w: w not in keep_stopwords, stopwords.words('russian')))

# определение того, подходит ли команда по н-граммам
def detect_match(lexemes: list[str], matchers: tuple[tuple[str]]) -> bool:
    for matcher in matchers:
        for ngram in ngrams(lexemes, len(matcher)):
            if ngram == matcher:
                return True
    return False

# Предполагается, что команда состоит из одного предложения
def parse_command(command_string: str):
    print('input:', command_string)
    tokens = word_tokenize(command_string)
    # print('tokens:', tokens)
    tokens = list(filter(lambda w: w not in stopwords, tokens))
    # print('filtered tokens:', tokens)
    lexemes = tokens # list(map(normalize_token, tokens))
    print('lexemes:', lexemes)
    for cmd, matchers in tuple(): # cmd_matchers_pairs:
        if detect_match(lexemes, matchers):
            return cmd()
    else:
        print('could not find any command')
        


class CommandServerHandler(server.BaseHTTPRequestHandler):
    def do_GET(self):
        try:
            query: dict = parse_qs(urlparse(self.path).query)
            print('query raw:', query)
            if 'query' not in query.keys():
                self.send_response(http.HTTPStatus.BAD_REQUEST, '"query" query not provided')
                return

            command = query['query'][0]
            print('query:', command)
            cmd = command # parse_command(command)

            self.send_response(http.HTTPStatus.OK)
            self.send_header('Content-Type', 'aplication/json')
            self.end_headers()
            self.wfile.write(bytes(cmd, 'utf-8'))

        except:
            self.send_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, 'failed to process request')

if __name__ == "__main__":
    # тест
    for cmd_string in '''\
    '''.splitlines():
        parse_command(cmd_string)
        print()

    host = os.getenv('HOST', '0.0.0.0')
    port = int(os.getenv('PORT', '8085'))

    commandServer = server.HTTPServer((host, port), CommandServerHandler)
    print('server running at', host, port)
    try:
        commandServer.serve_forever()
    except KeyboardInterrupt:
        pass

    commandServer.shutdown()
    print('server stopped')
