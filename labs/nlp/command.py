import pymorphy2
from abc import ABC, abstractmethod
from genreproc import get_genre, get_genres, genre_info
from user import UserContext
from animeapi import AnimeService


class Command(ABC):
    def __init__(self, tag):
        self.__tag = tag
        
    @property
    def tag(self):
        return self.__tag
    
    @abstractmethod
    def check(self, stems) -> float:
        pass
    
    @abstractmethod
    def apply(self, tokens) -> str:
        pass


def tags_check(tokens, tag_groups):
    def tag_check(tokens, tags):
        return len([t for t in tags if t in tokens]) / len(tags)
    return max(tag_check(tokens, grp) for grp in tag_groups)


class TotalInfoCommand(Command):
    def __init__(self):
        super().__init__('total-info')
        self.total_animes_in_db = 6669

    def check(self, stems) -> float:
        tags = (('сколько',), ('много', 'есть'))
        pun = 0
        if 'жанр' in stems:
            pun = 0.5
        return max(0, tags_check(stems, tags) - pun)

    def apply(self, tokens) -> str:
        return f'В базе хранится {self.total_animes_in_db} различных аниме фильмов и сериалов.'


class GenreInfoCommand(Command):
    def __init__(self, morph: pymorphy2.MorphAnalyzer):
        super().__init__('genre-info')
        self.__genre = None
        self.__morph = morph

    def check(self, stems) -> float:
        tags = (('какой', 'жанр'), ('много', 'жанр'))
        tag_val = tags_check(stems, tags)
        # check info about specific genre
        gsim, ggenre = get_genre(stems)
        return max(tag_val, 0.4 * gsim)

    def apply(self, tokens) -> str:
        stems = [self.__morph.normal_forms(t)[0] for t in tokens]
        tags = (('какой', 'жанр'), ('много', 'жанр'))
        tag_val = tags_check(stems, tags)
        gsim, ggenre = get_genre(tokens)
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

    def check(self, stems) -> float:
        # if there is no genre in message - not this command
        sim, genre = get_genre(stems)
        if sim < 0.6:
            return 0.0
        tags = (('любить',), ('нравиться',), ('понравиться',), ('хотеть',), ('хотеть', 'смотреть'), ('хотеть', 'посмотреть'))
        tag_val = tags_check(stems, tags)
        anti_tags = (('не', 'любить'), ('не', 'нравиться'), ('не', 'понравиться'), ('не', 'хотеть'), ('не', 'хотеть', 'смотреть'), ('не', 'хотеть', 'посмотреть'))
        anti_tag_val = tags_check(stems, anti_tags)
        return max(tag_val - 0.3 * anti_tag_val, 0)

    def apply(self, tokens) -> str:
        stems = [self.__morph.normal_forms(t)[0] for t in tokens]
        genres = get_genres(stems, 0.6)
        self.__user_context.like_genres(genres)
        
        recomendations = self.__service.get_recomendations()
        recomendations = [f'{i+1}) [ID#{a.id}] {a.title}'
                          for i, a in enumerate(recomendations)]
        recomendations = "\n".join(recomendations)
        res = f'Вам понравились следующие жанры: {", ".join(genres)}. Вот, что я могу Вам предложить:\n{recomendations}'
        print(res)
        return res


class GenreDislikeCommand(Command):
    def __init__(self, anime_service: AnimeService, morph: pymorphy2.MorphAnalyzer, user_context: UserContext):
        super().__init__('genre-dislike')
        self.__service = anime_service
        self.__morph = morph
        self.__user_context = user_context

    def check(self, stems) -> float:
        # if there is no genre in message - not this command
        sim, genre = get_genre(stems)
        if sim < 0.6:
            return 0.0
        tags = (('не', 'любить'), ('не', 'нравиться'), ('не', 'понравиться'), ('не', 'хотеть'), ('не', 'хотеть', 'смотреть'), ('не', 'хотеть', 'посмотреть'))
        tag_val = tags_check(stems, tags)
        return tag_val

    def apply(self, tokens) -> str:
        stems = [self.__morph.normal_forms(t)[0] for t in tokens]
        genres = get_genres(stems, 0.6)
        self.__user_context.dislike_genres(genres)
        
        recomendations = self.__service.get_recomendations()
        recomendations = [f'{i+1}) [ID#{a.id}] {a.title}'
                          for i, a in enumerate(recomendations)]
        recomendations = "\n".join(recomendations)
        return f'Вам не понравились следующие жанры: {", ".join(genres)}. Вот, что я могу Вам предложить:\n{recomendations}'
