from abc import ABC, abstractmethod
from animeapi import AnimeService, AnimeDTO
from genreproc import genres, genre_info
from grammar import ParseResult, QueryType
from pymorphy2 import MorphAnalyzer
from user import UserContext


class Command(ABC):
    def __init__(self, tag, queryType: QueryType):
        self.__tag = tag
        self.__query_type = queryType
        
    @property
    def tag(self):
        return self.__tag
    
    def check(self, query: ParseResult) -> bool:
        return self.__query_type == query.type
    
    @abstractmethod
    def apply(self, query: ParseResult) -> str:
        pass


class TotalInfoCommand(Command):
    def __init__(self):
        super().__init__('total-info', QueryType.GENERAL_QUESTION_1)
        self.total_animes_in_db = 6669

    def apply(self, query: ParseResult) -> str:
        return f'В базе хранится {self.total_animes_in_db} различных аниме фильмов и сериалов.'


class GenreTotalInfoCommand(Command):
    def __init__(self, morph: MorphAnalyzer):
        super().__init__('genre-total-info', QueryType.GENERAL_QUESTION_2)
        self.__morph = morph

    def apply(self, query: ParseResult) -> str:
        n = len(genres)
        gstr = self.__morph.parse('жанр')[0].make_agree_with_number(n)
        return f'В базе представлено {n} различных {gstr.word} начиная от психологического и заканчивая драмой.'


class GenreInfoCommand(Command):
    def __init__(self):
        super().__init__('genre-info', QueryType.GENERAL_QUESTION_3)
    
    def apply(self, query: ParseResult) -> str:
        return genre_info(query.genres[0][3])


class AnimeInfoCommand(Command):
    def __init__(self):
        super().__init__('anime-info', QueryType.GENERAL_QUESTION_4)
    
    def apply(self, query: ParseResult) -> str:
        a: AnimeDTO = query.anime_names[0][2]
        return f'Вот, что я могу рассказать про аниме {a.title}...'


class GenreLikeCommand(Command):
    def __init__(self, anime_service: AnimeService, user_context: UserContext):
        super().__init__('genre-like', QueryType.PREFERENCE_INFO_1)
        self.__service = anime_service
        self.__user_context = user_context

    def apply(self, query: ParseResult) -> str:
        genres = [g[3] for g in query.genres]
        self.__user_context.like_genres(genres)
        recs = self.__service.get_recomendations_str()
        return f'Вам понравились следующие жанры: {", ".join(genres)}. Вот, что я могу Вам предложить:\n{recs}'


class GenreDislikeCommand(Command):
    def __init__(self, anime_service: AnimeService, user_context: UserContext):
        super().__init__('genre-dislike', QueryType.PREFERENCE_INFO_2)
        self.__service = anime_service
        self.__user_context = user_context
    
    def apply(self, query: ParseResult) -> str:
        genres = [g[3] for g in query.genres]
        self.__user_context.dislike_genres(genres)
        recs = self.__service.get_recomendations_str()
        return f'Вам не понравились следующие жанры: {", ".join(genres)}. Вот, что я могу Вам предложить:\n{recs}'


class AnimeLikeCommand(Command):
    def __init__(self, anime_service: AnimeService, user_context: UserContext):
        super().__init__('anime-like', QueryType.PREFERENCE_INFO_3)
        self.__service = anime_service
        self.__user_context = user_context

    def apply(self, query: ParseResult) -> str:
        anime = query.anime_names[0][2]
        self.__user_context.like_anime(anime.id)
        recs = self.__service.get_recomendations_str()
        return f'Вам понравилось аниме {anime.title}. Вот, что я могу Вам предложить:\n{recs}'


class AnimeDislikeCommand(Command):
    def __init__(self, anime_service: AnimeService, user_context: UserContext):
        super().__init__('anime-dislike', QueryType.PREFERENCE_INFO_4)
        self.__service = anime_service
        self.__user_context = user_context

    def apply(self, query: ParseResult) -> str:
        anime = query.anime_names[0][2]
        self.__user_context.dislike_anime(anime.id)
        return f'Вычеркнул {anime.title} из списка рекомендаций.'
