import pymorphy2
from nltk.corpus import stopwords
from nltk.util import ngrams
from nltk.tokenize import word_tokenize
from user import UserContext
import command

# user contextual info about what he likes and dislikes
context = UserContext()

morph = pymorphy2.MorphAnalyzer()
keep_stopwords = ('чего', 'есть', 'много')
stopwords = list(filter(lambda w: w not in keep_stopwords, stopwords.words('russian')))

commands = [
    command.TotalInfoCommand(),
    command.GenreInfoCommand(morph),
    command.GenreLikeCommand(),
]

# Предполагается, что команда состоит из одного предложения
def parse_command(message: str) -> str:
    print('input:', message)
    tokens = word_tokenize(message)
    print('tokens:', tokens)
    tokens = list(filter(lambda w: w not in stopwords, tokens))
    print('filtered tokens:', tokens)
    stems = [morph.normal_forms(t) for t in tokens]
    print('stems:', stems)
    stems = [s[0] for s in stems]
    
    assurances = [(cmd.check(stems), cmd) for cmd in commands if cmd.check(stems) > 0]
    if len(assurances) == 0:
        print('not found command for message:', message)
        return 'Извините, я не смог распознать запрос'
    cmd = assurances[0][1]
    print('command tag:', cmd.tag)
    return cmd.apply(tokens)
