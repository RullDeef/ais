import os
import http
import traceback
from http import server
from urllib.parse import urlparse, parse_qs
import processor


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
            cmd = processor.parse_command(command)

            self.send_response(http.HTTPStatus.OK)
            self.send_header('Content-Type', 'text/plain')
            self.end_headers()
            self.wfile.write(bytes(cmd, 'utf-8'))

        except Exception as e:
            self.send_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, 'failed to process request')
            traceback.print_exc()


def run_server():
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


if __name__ == "__main__":
    import sys
    sys.stdout = sys.stderr
    
    run_server()
    exit()
    
        # Сколько аниме фильмов и сериалов хранится в базе данных?
        # Много чего есть?
        # Какие жанры представлены в базе?
        # Что из себя представляет жанр сейнен?
    # тест
    for cmd_string in '''\
Мне нравятся приключенческие аниме, также люблю посмотреть драму или комедию.
    '''.splitlines():
        print('result:', processor.parse_command(cmd_string))
        print()

