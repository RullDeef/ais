import animeapi
import command
import traceback
from grammar import SentenceParser
from morph import morph
from user import UserContext

api_server = animeapi.ApiServer('localhost', 8080)
anime_service = animeapi.AnimeService(api_server)
parser = SentenceParser(anime_service, morph)

# user contextual info about what he likes and dislikes
context = UserContext(anime_service)

commands = [
    command.TotalInfoCommand(),
    command.GenreTotalInfoCommand(morph),
    command.GenreInfoCommand(),
    command.AnimeInfoCommand(),
    command.GenreLikeCommand(anime_service, context),
    command.GenreDislikeCommand(anime_service, context),
    command.AnimeLikeCommand(anime_service, context),
    command.AnimeDislikeCommand(anime_service, context),
]

# Предполагается, что команда состоит из одного предложения
def parse_command(message: str) -> str:
    try:
        query = parser.parse(message)
        for cmd in commands:
            if not cmd.check(query):
                continue
            print('command tag:', cmd.tag)
            return cmd.apply(query)
    except Exception as e:
        print(traceback.format_exc())
        print(e)
        return 'Извините, я не смог распознать запрос.'
