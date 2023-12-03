import pymorphy2
from abc import ABC, abstractmethod
from typing import Optional
from strdist import str_similarity


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
    genres = (
        ("Psychological",),
		("Action",),
		("Shounen",),
		("Supernatural", "сверхъестественное"),
		("Sports",),
		("Martial Arts", "боевые искусства", "боевик"),
		("Historical", "исторический"),
		("Demons",),
		("Josei",),
		("Space", "космос"),
		("Mystery",),
		("Vampire",),
		("Cars", "машины", "тачки", "гонки"),
		("Super Power", "супер сила"),
		("Seinen",),
		("Sci-Fi", "сай-фай", "нучное", "научный", "научная фантастика"),
		("Magic", "магический", "волшебный"),
		("Parody",),
		("Thriller",),
		("Music",),
		("Game", "игра", "игровой"),
		("Fantasy",),
		("Adventure", "приключение"),
		("Romance",),
		("Police",),
		("Drama",),
		("Samurai",),
		("School", "школа"),
		("Comedy",),
		("Shoujo",),
		("Military", "военный"),
		("Horror", "ужастик"),
		("Slice of Life",),
		("Mecha", "робот"),
    )
    
    infos = {
        "Psychological": "Аниме в жанре психологическое (Psychological) обладает глубокими психологическими аспектами, ставящими персонажей в сложные ситуации, рассматривая их внутренние конфликты и психологические состояния.",
		"Action": "Экшон (Action) включает множество динамичных сцен схваток и битв, а также увлекательные и напряженные сюжеты.",
		"Shounen": "Шонен (Shounen) - жанр аниме, ориентированный на молодую аудиторию мужского пола. Часто включает в себя приключения, дружбу и боевые сражения.",
		"Supernatural": "Сверхъестественное (Supernatural). Аниме этого жанра добавляет в свой сюжет элементы сверхъестественных явлений, магии и необъяснимых событий.",
		"Sports": "Аниме о спорте фокусируется на соревнованиях, поднимаясь от развития спортсменов до их успехов и поражений.",
		"Martial Arts": "Боевые искусства (Martial Arts) - жанр аниме, где основное внимание уделяется боевым искусствам, сражениям и тренировкам.",
		"Historical": "Аниме исторического жанра (Historical) разворачивается в исторических сеттингах, воссоздавая эпохальные события и персонажей.",
		"Demons": "Аниме в жанре демоны (Demons) представляет собой миры, где демоны, ад и магия играют ключевую роль.",
		"Josei": "Жозей (Josei) - жанр, ориентированный на женскую аудиторию с более зрелыми темами и отношениями.",
		"Space": "Аниме в жанре космос (Space) разворачивается в космическом пространстве, включая космические оперы и приключения.",
		"Mystery": "Аниме в жанре мистики (Mystery) обладает загадочным сюжетом, зачастую включает в себя расследования, интриги и неожиданные повороты событий.",
		"Vampire": "Жанр аниме о вампирах (Vampire) представляет истории о существах ночи, мистических способностях и борьбе между светлым и темным.",
		"Cars": "Аниме о гонках (Cars) включает в себя захватывающие соревнования, скорость, адреналин и иногда технические новшества.",
		"Super Power": "Жанр аниме со суперсилами (Super Power) предлагает захватывающие битвы и приключения, где персонажи обладают уникальными способностями.",
		"Seinen": "Сейнен-аниме обращено к взрослой аудитории и включает в себя разнообразные темы, от реальной жизни до экшена и драмы.",
		"Sci-Fi": "Аниме научной фантастики (Sci-Fi) переносит нас в мир технологических новшеств, космических путешествий и футуристических обществ.",
		"Magic": "Жанр аниме с магией (Magic) предлагает волшебные приключения, волшебников, загадки и удивительные магические возможности.",
		"Parody": "Аниме в жанре пародии (Parody) предлагает забавные и ироничные адаптации знаменитых сюжетов и анимационных стилей.",
		"Thriller": "Жанр триллера в аниме представляет напряженные и драматические сюжеты с неожиданными поворотами и интригой.",
		"Music": "Аниме о музыке фокусируется на страсти к музыке, музыкальных группах и творческом процессе.",
		"Game": "Жанр аниме об играх включает в себя захватывающие сюжеты об игровых мирах, сражениях и увлекательных приключениях.",
		"Fantasy": "Аниме в жанре фэнтэзи создает миры с драконами, эльфами, магией и приключениями по удивительному и волшебному фэнтэзийному миру.",
		"Adventure": "Жанр приключений представляет захватывающие истории, путешествия и открытия в удивительных мирах.",
		"Romance": "Аниме о романтике предлагает теплые и трогательные истории любви, сложных отношений и моментов открытия сердец.",
		"Police": "Жанр аниме о полиции отображает расследования, детективные истории, правосудие и борьбу с преступностью.",
		"Drama": "Аниме драмы погружает зрителя в интенсивные и эмоциональные истории, о различных аспектах жизни и взаимоотношениях.",
		"Samurai": "Жанр аниме о самураях представляет культуру, боевые навыки и душевные качества воинов-самураев.",
		"School": "Аниме о школе предлагает истории о школьной жизни, дружбе, любви и жизненных уроках.",
		"Comedy": "Жанр комедии в аниме приносит улыбки и смех, создавая забавные и занимательные истории.",
		"Shoujo": "Шоджо (Shoujo) - жанр аниме, ориентированный на женскую аудиторию с яркими романтическими и драматическими сюжетами.",
		"Military": "Военное (Military) - жанр аниме, посвященный вооруженным силам, сражениям и военным стратегиям.",
		"Horror": "Аниме ужасов предлагает атмосферу напряжения, мистики и демонстрирует умение пугать.",
		"Slice of Life": "Жизнь во всей её полноте (Slice of Life) - жанр аниме, в котором описываются повседневные истории, реалистичные персонажи и человеческие взаимоотношения.",
		"Mecha": "Аниме в жарне меха (Mecha) представляет боевые и гигантские роботы, механические сражения и научно-фантастические истории.",
    }
    
    def __init__(self, morph: pymorphy2.MorphAnalyzer):
        super().__init__('genre-info')
        self.__genre = None
        self.__morph = morph

    def check(self, stems) -> float:
        tags = (('какой', 'жанр'), ('много', 'жанр'))
        tag_val = tags_check(stems, tags)
        # check info about specific genre
        gsim, ggenre = 0.0, ''
        ginfo = self.__extract_genre(stems)
        if ginfo is not None:
            gsim, ggenre = ginfo
        return max(tag_val, gsim)
    
    def __extract_genre(self, stems: list[str]) -> Optional[tuple[float, str]]:
        if 'жанр' not in stems:
            return None
        pos = stems.index('жанр')
        if pos == len(stems) - 1:
            return None
        genre = stems[pos + 1]
        sim_map = [
            (max(str_similarity(genre, gg) for gg in g), g[0])
            for g in GenreInfoCommand.genres]
        return max(sim_map, key=lambda s: s[0])

    def apply(self, tokens) -> str:
        stems = [self.__morph.normal_forms(t)[0] for t in tokens]
        tags = (('какой', 'жанр'), ('много', 'жанр'))
        tag_val = tags_check(stems, tags)
        ginfo = self.__extract_genre(stems)
        if ginfo is None or tag_val > 0.6:
            n = len(GenreInfoCommand.genres)
            gstr = self.__morph.parse('жанр')[0].make_agree_with_number(n)
            return f'В базе представлено {n} различных {gstr.word} начиная от психологического и заканчивая драмой.'
        sim, genre = ginfo
        if sim < 0.6:
            return f'Извините, я не понял какой жанр Вы имели в виду. Возможно {genre}?'
        return GenreInfoCommand.infos[genre]


class GenreLikeCommand(Command):
    def __init__(self):
        super().__init__('genre-like')

    def check(self, stems) -> float:
        # TODO
        return 0.0

    def apply(self, tokens) -> str:
        # TODO
        return ''
