from animeapi import AnimeService
from durproc import is_duration
from genreproc import is_genre
from nltk.grammar import CFG
from nltk.corpus import stopwords as nltk_stopwords
from nltk.parse import ChartParser
from nltk.tokenize import word_tokenize
from nltk.util import ngrams
from pymorphy2 import MorphAnalyzer

grammar = CFG.fromstring("""
    # грамматика запроса
    S -> GQ | PR | FR
    
    # общие вопросы
    GQ -> GQ1 | GQ2 | GQ3 | GQ4
    
    # запросты-предпочтения
    PR -> PR1 | PR2 | PR3 | PR4
    
    # запросы-фильтры
    FR -> FR1 | FR2
    
    # Общий вопрос №1 -----------------------------------------------
    # вопрос про количество аниме в базе
    GQ1 -> HowMany 'чего' 'есть' OptionalQuestionMark
    GQ1 -> HowMany AnimeNounSP OptionalQuestionMark
    GQ1 -> HowMany AnimeNounSP ContainsIn Database OptionalQuestionMark
    
    # фраза, означающая вопросительное "сколько?"
    HowMany -> 'сколько' | 'много' | 'много' 'ли'
    
    # фраза, означающая одно или несколько аниме
    AnimeNounSP -> AnimeNounS | AnimeNounP
    AnimeNounS -> 'аниме' | 'аниме' 'фильм' | 'аниме' 'сериал'
    AnimeNounP -> 'аниме' 'фильмы' | 'аниме' 'сериалы'
    
    # фраза, означающая содержание чего-то в чем-то
    ContainsIn -> Contains | Contains 'в' | 'в'
    Contains -> 'храниться' | 'содержаться' | 'находиться' | 'представить'
    
    # База данных
    Database -> 'база' | 'база' 'данные' | 'хранилище' | 'память'

    # Общий вопрос №2 -----------------------------------------------
    # вопрос про количество жанров в базе
    GQ2 -> HowMany GenreNoun OptionalQuestionMark
    GQ2 -> HowMany GenreNoun ContainsIn Database OptionalQuestionMark
    GQ2 -> 'какой' GenreNoun ContainsIn Database OptionalQuestionMark
    
    GenreNoun -> 'жанр' | 'жанры'
    
    # Общий вопрос №3 -----------------------------------------------
    # вопрос про смысл конкретного жанра
    GQ3 -> ConcreteGenre Is What OptionalQuestionMark
    GQ3 -> What CanTell About ConcreteGenre OptionalQuestionMark
    GQ3 -> 'что' 'представлять' ConcreteGenre OptionalQuestionMark
    
    Is -> '-' | 'это' | '-' 'это'
    What -> 'что' | 'что' 'такой' | 'ху'
    
    # продукции для определения конкретных жанров
    ConcreteGenre -> 'жанр' '#GENRE#' | '#GENRE#'
    
    # Общий вопрос №4 -----------------------------------------------
    # получение информации по конкретному аниме
    GQ4 -> What CanTell About ConcreteAnime OptionalQuestionMark
    
    CanTell -> 'мочь' 'рассказать'
    About -> 'о' | 'обо' | 'про' |
    ConcreteAnime -> '#ANIME#' | AnimeNounS '#ANIME#'
    
    # Запрос-предпочтения №1 ----------------------------------------
    # нравятся жанры
    PR1 -> ILike ConcreteGenreP
    
    ILike -> 'я' Like | Like
    Like -> 'любить' | 'нравиться' | 'понравиться'
    
    ConcreteGenreP -> ConcreteGenre | ConcreteGenre ConcreteGenreP
    ConcreteGenreP -> ConcreteGenre 'и' ConcreteGenreP
    ConcreteGenreP -> ConcreteGenre ',' ConcreteGenreP

    # Запрос-предпочтения №2 ----------------------------------------
    # не нравятся жанры
    PR2 -> IDislike ConcreteGenreP
    
    IDislike -> 'я' Dislike | Dislike
    Dislike -> 'не' Like
    
    # Запрос-предпочтения №3 ----------------------------------------
    # нравится аниме
    PR3 -> ILike ConcreteAnime
    
    # Запрос-предпочтения №4 ----------------------------------------
    # не нравится аниме
    PR4 -> IDislike ConcreteAnime
    
    # Запрос-фильтр №1 ----------------------------------------------
    # отфильтровать по длительности
    FR1 -> FilterKeep ConcreteDurationP
    FR1 -> FilterKeep ConcreteDurationP AnimeNounSP
    
    FilterKeep -> FilterKeepVerb | FilterKeepVerb 'только'
    FilterKeepVerb -> 'отфильтровать' | 'оставить'
    
    # Запрос-фильтр №2 ----------------------------------------------
    # отфильтровать по длительности
    FR2 -> FilterDrop ConcreteDurationP
    FR2 -> FilterDrop ConcreteDurationP AnimeNounSP
    
    FilterDrop -> 'удалить' | 'убрать'
    
    ConcreteDurationP -> ConcreteDuration
    ConcreteDurationP -> ConcreteDuration ConcreteDurationP
    ConcreteDurationP -> ConcreteDuration ',' ConcreteDurationP
    ConcreteDurationP -> ConcreteDuration 'и' ConcreteDurationP
    ConcreteDuration -> '#DUR#'

    # опциональный вопросительный знак
    OptionalQuestionMark -> '?' |
""")


class QueryType:
    GENERAL_QUESTION_1 = 'GQ1'
    GENERAL_QUESTION_2 = 'GQ2'
    GENERAL_QUESTION_3 = 'GQ3'
    GENERAL_QUESTION_4 = 'GQ4'
    PREFERENCE_INFO_1 = 'PR1'
    PREFERENCE_INFO_2 = 'PR2'
    PREFERENCE_INFO_3 = 'PR3'
    PREFERENCE_INFO_4 = 'PR4'
    FILTER_REQUEST_1 = 'FR1'
    FILTER_REQUEST_2 = 'FR2'


class ParseResult:
    def __init__(self, sentence: str, type: QueryType, genres: list, durations: list, anime_names: list):
        self.sentence = sentence
        self.type = type
        self.genres = genres
        self.durations = durations
        self.anime_names = anime_names
    
    def __repr__(self) -> str:
        return f'ParseResult("{self.sentence}", {self.type}, {self.genres}, {self.durations}, {self.anime_names})'


class SentenceParser:
    def __init__(self, anime_service: AnimeService, morph: MorphAnalyzer):
        self.__morph = morph
        self.__anime_service = anime_service
        self.__parser = ChartParser(grammar)
        self.__init_stopwords()
    
    def __init_stopwords(self):
        keep_stopwords = ('чего', 'есть', 'много', 'не')
        stopwords = list(filter(lambda w: w not in keep_stopwords, nltk_stopwords.words('russian')))
        stopwords += ['также']
        self.__stopwords = stopwords
    
    def parse(self, sentence: str):
        print('SentenceParser| sentence:', sentence)
        tokens = word_tokenize(sentence)
        print('SentenceParser| tokens:', tokens)
        tokens = list(filter(lambda w: w not in self.__stopwords, tokens))
        print('SentenceParser| filtered tokens:', tokens)
        replaced_animes = self.__extract_anime_names(tokens)
        if len(replaced_animes) > 0:
            print('Sentenceparser| animes:', replaced_animes)
        replaced_genres = self.__extract_genres(tokens)
        if len(replaced_genres) > 0:
            print('SentenceParser| genres:', replaced_genres)
        replaced_durations = self.__extract_durations(tokens)
        if len(replaced_durations) > 0:
            print('SentenceParser| durations:', replaced_durations)
        stems = [[t] if '#' in t else self.__morph.normal_forms(t) for t in tokens ]
        print('SentenceParser| all stems:', stems)
        stems = [s[0] for s in stems]
        print('SentenceParser| selected stems:', stems)
        tree = self.__parser.parse_one(stems)
        print('parsed tree:', tree)
        type = tree[0][0].label()
        res = ParseResult(sentence, type, replaced_genres, replaced_durations, replaced_animes)
        print('SentenceParser| res:', res)
        return res
        
    def __extract_anime_names(self, tokens: list[str]) -> list:
        if 'аниме' not in tokens:
            return []
        animes = []
        for k in (3, 2, 1):
            for i, ngram in enumerate(ngrams(tokens, k)):
                if 'аниме' in ngram:
                    continue
                name = ' '.join(ngram)
                anime = self.__anime_service.find_anime_exact(name)
                if anime is not None:
                    tokens[i:i+k] = ['#ANIME#']
                    animes.append(((i, k), name, anime))
        return animes

    def __extract_genres(self, tokens: list[str]) -> list:
        replaced_genres = []
        for i, t in enumerate(tokens):
            prob, gen = is_genre(t)
            if prob > 0.7:
                tokens[i] = '#GENRE#'
                replaced_genres.append((i, prob, t, gen))
        return replaced_genres

    def __extract_durations(self, tokens: list[str]) -> list:
        replaced_durations = []
        for k in (3, 2, 1):
            for i, ngram in enumerate(ngrams(tokens, k)):
                prob, dur = is_duration(' '.join(ngram))
                if prob != 0:
                    tokens[i:i+k] = ['#DUR#']
                    replaced_durations.append(((i, k), prob, dur))
        return replaced_durations


if __name__ == "__main__":
    from morph import morph
    from animeapi import ApiServer
    
    api_server = ApiServer('localhost', '8080')
    anime_service = AnimeService(api_server)
    
    parser = SentenceParser(anime_service, morph)
    
    res = parser.parse("Сколько жанров представлено в базе данных?")
    print(res)
