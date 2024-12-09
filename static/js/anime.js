// code for dark-mode
const htmlElement = document.documentElement;
const themeSwitch = document.querySelector('.switch input');

// Initialize theme based on localStorage
document.addEventListener('DOMContentLoaded', () => {
    const darkMode = localStorage.getItem('darkMode') === 'true';
    themeSwitch.checked = darkMode;
    if (darkMode) {
        htmlElement.classList.add('dark-mode');
    } else {
        htmlElement.classList.remove('dark-mode');
    }
});

//Toggle dark mode on checkbox change
themeSwitch.addEventListener('change', () => {
    if (themeSwitch.checked) {
        htmlElement.classList.add('dark-mode');
        localStorage.setItem('darkMode', 'true');
    } else {
        htmlElement.classList.remove('dark-mode');
        localStorage.setItem('darkMode', 'false');
    }
});

// Creds
const accessToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6ImE3ZjI5OTg3NjUxMmJjNThiZGU5MGFkYzllN2U0NzZmYzhhYWU2OWE5MWNlZTI5NDVlMDMwNGE2NTY0MzVhNzQwODUwODIzMzg5YjdiOGZjIn0.eyJhdWQiOiIyMTc4MyIsImp0aSI6ImE3ZjI5OTg3NjUxMmJjNThiZGU5MGFkYzllN2U0NzZmYzhhYWU2OWE5MWNlZTI5NDVlMDMwNGE2NTY0MzVhNzQwODUwODIzMzg5YjdiOGZjIiwiaWF0IjoxNzMyMTg0NDY3LCJuYmYiOjE3MzIxODQ0NjcsImV4cCI6MTc2MzcyMDQ2Nywic3ViIjoiNTgwOTAwNCIsInNjb3BlcyI6W119.sLQd_ULrAtUNwGiwrGPSbAT4cnvUtvpc26D3k_UYRxPJ9JbdvJcLPB2DTSF-u_YoY5fuMjEuejs7-iyGprZWFHfHHARnvY757wj8V9PT9yPVXFiUWZUnm9ApagSEEAoW2PqpNiu1Y6BWYELHatF3rd5_78ZklIL_GC4PKxdo-kcLc6J03B7SJoOpq64D9yyOZClvy_AWVeLDEVK0RH8h583n16jmRwigD50MBUqIfWgYtz1Lkb3Uv8_STpdB6v0ngy8C4GoakExXB8ieaOnxylx94ZkuJ2Hzyy-vuY4yrzIplrEkOcj1nyXKKlOyU4cnn4K4O9aFeNrCd_8rON0B36alNEARONXwVCngODbtMsBz0Hes09PgxBsq-sJU2YVAM8ObAPEtgo0h6SxtoD5sFyHH_RiaK9nhegZMGDGprzP5m8yNQCrd7yKjJqiUMz9yfEDRcPZAABFCc9r2Gdft1uIpvD4nBtU4MqRm7oH1vLCqshjIIT9MKpCBT_JbdSxovyjrfB50GtbqMrCw1rTTViGWwIsHuCAj2WD-UjMZS3zOz5uaE4J2DzMTTILYF4dOHMxxs12do1fHRv4qseazNLnidyrRbET1GvUyRQ_mdvDZuXROcVgm7uwU22gT_Nwh5IDV7c3yqkUs27hwoltMQ88Krd47EsfGlrCrjA0JwXs";

async function fetchPopularAnime() {
    const query = `
        query {
            Page(page: 1, perPage: 10) {
                media(type: ANIME, sort: POPULARITY_DESC) {
                    title {
                        romaji
                        english
                    }
                    description    
                    bannerImage
                }
            }
        }
    `;

    const headers = {
        'Authorization': 'Bearer ' + accessToken,
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    };

    const response = await fetch('https://graphql.anilist.co', {
        method: 'POST',
        headers: headers,
        body: JSON.stringify({ query: query }),
    });

    if (!response.ok) {
        throw new Error('Failed to fetch data from Anilist API')
    }
    const data = await response.json();
    return data.data.Page.media;
}

async function populateSlideShow() {
    const slideShow = document.getElementById('slideshow');
    const series = await fetchPopularAnime();
    series.forEach((tv, index) => {
        const slide = document.createElement('div');
        slide.classList.add('slide');
        if (index === 0) slide.classList.add('active');

        slide.innerHTML = `
        <img src="${tv.bannerImage}" alt="${tv.title.romaji}">
        <div class="slide-content">
            <h2>${tv.title.english || tv.title.romaji}</h2>
            <p>${tv.description}"</p>
            <button></button>
            <button>Info</button>
        </div>
        `;
        slideShow.appendChild(slide);
    });
    initializeSlideShow();
}

function initializeSlideShow() {
    const slides = document.querySelectorAll('.slide');
    let currentSlide = 0;

    const showSlide = (index) => {
        slides.forEach((slide, i) => {
            slide.classList.toggle('active', i === index);
        });
    };

    document.getElementById('prev-btn').addEventListener('click', () => {
        currentSlide = (currentSlide - 1 + slides.length) % slides.length;
        showSlide(currentSlide);
    });

    document.getElementById('next-btn').addEventListener('click', () => {
        currentSlide = (currentSlide + 1) % slides.length;
        showSlide(currentSlide);
    });
}

async function fetchAnimeByCategory(category, containerId) {
    const query = `
        query {
            Page(page: 1, perPage: 25) {
                media(type: ANIME, sort: ${category}) {
                    title {
                        romaji
                        english
                    }
                    coverImage {
                        large
                    }
                }
            }
        }
    `;

    const headers = {
        'Authorization': 'Bearer ' + accessToken,
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    };

    const response = await fetch('https://graphql.anilist.co', {
        method: 'POST',
        headers: headers,
        body: JSON.stringify({ query: query })
    });
    const data = await response.json();
    const container = document.getElementById(containerId);

    data.data.Page.media.forEach((anime) => {
        const animeCard = document.createElement('div');
        animeCard.classList.add('movie-card');
        animeCard.innerHTML = `
            <img src="${anime.coverImage.large}" alt="${anime.title.romaji}">
            <p>${anime.title.english || anime.title.romaji}</p>
        `;
        container.appendChild(animeCard);
    });
}

async function fetchReleasingAnime(containerId) {
    const query = `
        query {
            Page(page:1, perPage: 25) {
                media(type: ANIME, status: RELEASING, sort: POPULARITY_DESC) {
                    title {
                        romaji
                        english
                    }
                    coverImage {
                        medium
                        large
                    }
                }
            }
        }
    `;

    const headers = {
        'Authorization': 'Bearer ' + accessToken,
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    };

    try {
        const response = await fetch('https://graphql.anilist.co', {
            method: 'POST',
            headers: headers,
            body: JSON.stringify({ query: query }),
        });

        if (!response.ok) {
            throw new Error('Failed to fetch data from Anilist API');
        }

        const data = await response.json();
        const animeList = data.data.Page.media;
        const container = document.getElementById(containerId);

        animeList.forEach((anime) => {
            const animeCard = document.createElement('div');
            animeCard.classList.add('movie-card');
            animeCard.innerHTML = `
                <img src="${anime.coverImage.large}" alt="${anime.title.romanji}">
                <p>${anime.title.english || anime.title.romanji}</p>
            `;
            container.appendChild(animeCard);
        });
    } catch (error) {
        console.error('Error fetching popular anime:', error);
    }
}

async function initializeContent() {
    await populateSlideShow();
    await fetchAnimeByCategory('POPULARITY_DESC', 'popular-category');
    await fetchAnimeByCategory('TRENDING_DESC', 'trending-category');
    await fetchAnimeByCategory('SCORE_DESC', 'top-rated-category');
    await fetchReleasingAnime('currently-airing-category');
}

initializeContent();
