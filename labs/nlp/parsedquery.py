from nltk.corpus import stopwords
from nltk.util import ngrams
from nltk.tokenize import word_tokenize
from morph import morph

keep_stopwords = ('чего', 'есть', 'много', 'не')
stopwords = list(filter(lambda w: w not in keep_stopwords, stopwords.words('russian')))

stopwords += ['также']


class ParsedQuery:
    def __init__(self, query: str, tokens: list[str], stems: list[str]):
        self.__query = query
        self.__tokens = tokens
        self.__stems = stems

    @property
    def query(self) -> str:
        return self.__query
    
    @property
    def tokens(self) -> list[str]:
        return self.__tokens[:]

    @property
    def stems(self) -> list[str]:
        return self.__stems[:]

    @classmethod
    def parse(cls, query: str) -> 'ParsedQuery':
        print('input:', query)
        tokens = word_tokenize(query)
        print('tokens:', tokens)
        tokens = list(filter(lambda w: w not in stopwords, tokens))
        print('filtered tokens:', tokens)
        stems = [morph.normal_forms(t) for t in tokens]
        print('stems:', stems)
        stems = [s[0] for s in stems]
        return cls(query, tokens, stems)
