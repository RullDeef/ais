import animeapi
import command
from user import UserContext
from parsedquery import ParsedQuery
from morph import morph

api_server = animeapi.ApiServer('localhost', 8080)
anime_service = animeapi.AnimeService(api_server)

# user contextual info about what he likes and dislikes
context = UserContext(anime_service)

commands = [
    command.TotalInfoCommand(),
    command.GenreInfoCommand(morph),
    command.GenreLikeCommand(anime_service, morph, context),
    command.GenreDislikeCommand(anime_service, morph, context),
    command.AnimeLikeCommand(anime_service, context),
]

# Предполагается, что команда состоит из одного предложения
def parse_command(message: str) -> str:
    query = ParsedQuery.parse(message)
    
    assurances = sorted([(cmd.check(query), cmd) for cmd in commands], key=lambda a: -a[0])
    print('assurances:', [(round(score, 2), cmd.tag) for score, cmd in assurances])
    assurances = [a for a in assurances if a[0] >= 0.5]
    if len(assurances) == 0:
        print('not found command for message:', message)
        return 'Извините, я не смог распознать запрос.'
    cmd = assurances[0][1]
    print('command tag:', cmd.tag)
    return cmd.apply(query)
