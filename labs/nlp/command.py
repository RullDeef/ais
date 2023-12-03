import pymorphy2
from abc import ABC, abstractmethod
from typing import Optional
from genreproc import get_genre, genre_info


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
    def __init__(self):
        super().__init__('genre-like')

    def check(self, stems) -> float:
        # TODO
        return 0.0

    def apply(self, tokens) -> str:
        # TODO
        return ''
