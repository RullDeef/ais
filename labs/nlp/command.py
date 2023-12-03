import pymorphy2
from nltk.util import ngrams
from abc import ABC, abstractmethod
from genreproc import get_genre, get_genres, genre_info
from user import UserContext
from animeapi import AnimeService, AnimeDTO
from parsedquery import ParsedQuery
from strdist import str_similarity


class Command(ABC):
    def __init__(self, tag):
        self.__tag = tag
        
    @property
    def tag(self):
        return self.__tag
    
    @abstractmethod
    def check(self, query: ParsedQuery) -> float:
        pass
    
    @abstractmethod
    def apply(self, query: ParsedQuery) -> str:
        pass


def tags_check(tokens: list[str], tag_groups: list[tuple[str]]):
    def tag_check(tokens, tags):
        return len([t for t in tags if t in tokens]) / len(tags)
    return max(tag_check(tokens, grp) for grp in tag_groups)


class TotalInfoCommand(Command):
    def __init__(self):
        super().__init__('total-info')
        self.total_animes_in_db = 6669

    def check(self, query: ParsedQuery) -> float:
        tags = (('сколько',), ('много', 'есть'))
        pun = 0
        if 'жанр' in query.stems:
            pun = 0.5
        return max(0, tags_check(query.stems, tags) - pun)

    def apply(self, query: ParsedQuery) -> str:
        return f'В базе хранится {self.total_animes_in_db} различных аниме фильмов и сериалов.'


class GenreInfoCommand(Command):
    def __init__(self, morph: pymorphy2.MorphAnalyzer):
        super().__init__('genre-info')
        self.__morph = morph

    def check(self, query: ParsedQuery) -> float:
        tags = (('какой', 'жанр'), ('много', 'жанр'))
        tag_val = tags_check(query.stems, tags)
        # check info about specific genre
        gsim, ggenre = get_genre(query.tokens)
        return max(tag_val, 0.4 * gsim)

    def apply(self, query: ParsedQuery) -> str:
        tags = (('какой', 'жанр'), ('много', 'жанр'))
        tag_val = tags_check(query.stems, tags)
        gsim, ggenre = get_genre(query.tokens)
        if tag_val > 0.6:
            n = len(GenreInfoCommand.genres)
            gstr = self.__morph.parse('жанр')[0].make_agree_with_number(n)
            return f'В базе представлено {n} различных {gstr.word} начиная от психологического и заканчивая драмой.'
        if gsim < 0.6:
            return f'Извините, я не понял какой жанр Вы имели в виду. Возможно {ggenre}?'
        return genre_info(ggenre)


class GenreLikeCommand(Command):
    def __init__(self, anime_service: AnimeService, morph: pymorphy2.MorphAnalyzer, user_context: UserContext):
        super().__init__('genre-like')
        self.__service = anime_service
        self.__morph = morph
        self.__user_context = user_context

    def check(self, query: ParsedQuery) -> float:
        # if there is no genre in message - not this command
        sim, genre = get_genre(query.stems)
        if sim < 0.6:
            return 0.0
        tags = (('любить',), ('нравиться',), ('понравиться',), ('хотеть',), ('хотеть', 'смотреть'), ('хотеть', 'посмотреть'))
        tag_val = tags_check(query.stems, tags)
        anti_tags = (('не', 'любить'), ('не', 'нравиться'), ('не', 'понравиться'), ('не', 'хотеть'), ('не', 'хотеть', 'смотреть'), ('не', 'хотеть', 'посмотреть'))
        anti_tag_val = tags_check(query.stems, anti_tags)
        return max(tag_val - 0.3 * anti_tag_val, 0)

    def apply(self, query: ParsedQuery) -> str:
        genres = get_genres(query.stems, 0.6)
        self.__user_context.like_genres(genres)
        
        recs = self.__service.get_recomendations_str()
        return f'Вам понравились следующие жанры: {", ".join(genres)}. Вот, что я могу Вам предложить:\n{recs}'


class GenreDislikeCommand(Command):
    def __init__(self, anime_service: AnimeService, morph: pymorphy2.MorphAnalyzer, user_context: UserContext):
        super().__init__('genre-dislike')
        self.__service = anime_service
        self.__morph = morph
        self.__user_context = user_context

    def check(self, query: ParsedQuery) -> float:
        # if there is no genre in message - not this command
        sim, genre = get_genre(query.stems)
        if sim < 0.6:
            return 0.0
        tags = (('не', 'любить'), ('не', 'нравиться'), ('не', 'понравиться'), ('не', 'хотеть'), ('не', 'хотеть', 'смотреть'), ('не', 'хотеть', 'посмотреть'))
        tag_val = tags_check(query.stems, tags)
        return (tag_val + sim) / 2

    def apply(self, query: ParsedQuery) -> str:
        genres = get_genres(query.stems, 0.6)
        self.__user_context.dislike_genres(genres)
        
        recs = self.__service.get_recomendations_str()
        return f'Вам не понравились следующие жанры: {", ".join(genres)}. Вот, что я могу Вам предложить:\n{recs}'


def extract_anime(service: AnimeService, tokens: list[str]) -> tuple[float, AnimeDTO]:
    # remove all after аниме
    try:
        ind = tokens.index('аниме')
        tokens = tokens[ind+1:]
    except:
        return 0, ''

    # exact search
    title = " ".join(tokens)
    anime = service.find_anime_exact(title)
    if anime is not None:
        return 1.0, anime
    # fuzzy search
    max_anime, max_assurance = '', 0.0
    for n in range(len(tokens), 0, -1):
        for ngram in ngrams(tokens, n):
            title = ' '.join(ngram)
            anime = service.find_anime_fuzzy(title)[0]
            assurance = str_similarity(title, anime.title)
            if assurance > max_assurance:
                max_anime = anime
                max_assurance = assurance
    return max_assurance, max_anime


class AnimeLikeCommand(Command):
    def __init__(self, anime_service: AnimeService, user_context: UserContext):
        super().__init__('anime-like')
        self.__service = anime_service
        self.__user_context = user_context
        self.__anime = None
    
    def check(self, query: ParsedQuery) -> float:
        anti_tags = (('не', 'любить'), ('не', 'нравиться'), ('не', 'понравиться'), ('не', 'хотеть'), ('не', 'хотеть', 'смотреть'), ('не', 'хотеть', 'посмотреть'))
        anti_tag_val = tags_check(query.stems, anti_tags)
        tags = (('любить',), ('нравиться',), ('понравиться',), ('хотеть',), ('хотеть', 'смотреть'), ('хотеть', 'посмотреть'))
        tag_val = tags_check(query.stems, tags)
        tag_val = max(0, tag_val - 0.3 * anti_tag_val)
        if tag_val < 0.7:
            return 0.0
        assurance, anime = extract_anime(self.__service, query.tokens)
        if assurance < 0.7:
            return 0.0
        self.__anime = anime
        return max(0, assurance - 0.3 * anti_tag_val)

    def apply(self, query: ParsedQuery) -> str:
        self.__service.like_anime(self.__anime.id)
        recs = self.__service.get_recomendations_str()
        return f'Вам понравилось аниме {self.__anime.title}. Вот, что я могу Вам предложить:\n{recs}'


class AnimeDislikeCommand(Command):
    def __init__(self, anime_service: AnimeService, user_context: UserContext):
        super().__init__('anime-dislike')
        self.__service = anime_service
        self.__user_context = user_context
        self.__anime = None
    
    def check(self, query: ParsedQuery) -> float:
        tags = (('не', 'любить'), ('не', 'нравиться'), ('не', 'понравиться'), ('не', 'хотеть'), ('не', 'хотеть', 'смотреть'), ('не', 'хотеть', 'посмотреть'))
        tag_val = tags_check(query.stems, tags)
        if tag_val < 0.7:
            return 0.0
        assurance, anime = extract_anime(self.__service, query.tokens)
        if assurance < 0.7:
            return 0.0
        self.__anime = anime
        return assurance

    def apply(self, query: ParsedQuery) -> str:
        self.__service.dislike_anime(self.__anime.id)
        return f'Вычеркнул {self.__anime.title} из списка рекомендаций.'
