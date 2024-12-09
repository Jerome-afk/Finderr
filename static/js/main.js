/* ............................................... Fetch Media Content .................................................. */
document.addEventListener("DOMContentLoaded", () => {
fetchPopularMovies();
fetchPopularTVShows();
const accessToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6ImE3ZjI5OTg3NjUxMmJjNThiZGU5MGFkYzllN2U0NzZmYzhhYWU2OWE5MWNlZTI5NDVlMDMwNGE2NTY0MzVhNzQwODUwODIzMzg5YjdiOGZjIn0.eyJhdWQiOiIyMTc4MyIsImp0aSI6ImE3ZjI5OTg3NjUxMmJjNThiZGU5MGFkYzllN2U0NzZmYzhhYWU2OWE5MWNlZTI5NDVlMDMwNGE2NTY0MzVhNzQwODUwODIzMzg5YjdiOGZjIiwiaWF0IjoxNzMyMTg0NDY3LCJuYmYiOjE3MzIxODQ0NjcsImV4cCI6MTc2MzcyMDQ2Nywic3ViIjoiNTgwOTAwNCIsInNjb3BlcyI6W119.sLQd_ULrAtUNwGiwrGPSbAT4cnvUtvpc26D3k_UYRxPJ9JbdvJcLPB2DTSF-u_YoY5fuMjEuejs7-iyGprZWFHfHHARnvY757wj8V9PT9yPVXFiUWZUnm9ApagSEEAoW2PqpNiu1Y6BWYELHatF3rd5_78ZklIL_GC4PKxdo-kcLc6J03B7SJoOpq64D9yyOZClvy_AWVeLDEVK0RH8h583n16jmRwigD50MBUqIfWgYtz1Lkb3Uv8_STpdB6v0ngy8C4GoakExXB8ieaOnxylx94ZkuJ2Hzyy-vuY4yrzIplrEkOcj1nyXKKlOyU4cnn4K4O9aFeNrCd_8rON0B36alNEARONXwVCngODbtMsBz0Hes09PgxBsq-sJU2YVAM8ObAPEtgo0h6SxtoD5sFyHH_RiaK9nhegZMGDGprzP5m8yNQCrd7yKjJqiUMz9yfEDRcPZAABFCc9r2Gdft1uIpvD4nBtU4MqRm7oH1vLCqshjIIT9MKpCBT_JbdSxovyjrfB50GtbqMrCw1rTTViGWwIsHuCAj2WD-UjMZS3zOz5uaE4J2DzMTTILYF4dOHMxxs12do1fHRv4qseazNLnidyrRbET1GvUyRQ_mdvDZuXROcVgm7uwU22gT_Nwh5IDV7c3yqkUs27hwoltMQ88Krd47EsfGlrCrjA0JwXs"
fetchPopularAnime(accessToken);

if(localStorage.getItem('dark-mode') === 'true') {
    document.body.classList.add('dark-mode')
}
});

document.addEventListener('DOMContentLoaded', () => {
    const upperRow = document.getElementById('upper-row');
    const lowerRow = document.getElementById('lower-row');

    // Function to rearrange buttons dynamically
    function rearrangeButtons() {
        const allItems = [...upperRow.children, ...lowerRow.children];
        const maxPerRow = window.innerWidth <= 768 ? 2 : 4;

        // Clear rows
        upperRow.innerHTML = '';
        lowerRow.innerHTML = '';

        // Add items dynamically
        allItems.forEach((item, index) => {
            if (index < maxPerRow) {
                upperRow.appendChild(item);
            } else {
                lowerRow.appendChild(item);
            }
        });
    }

    // Initial arrangement
    rearrangeButtons();

    // Rearrange on window resize
    window.addEventListener('resize', rearrangeButtons);
});

function fetchPopularMovies() {
    fetch('https://api.themoviedb.org/3/movie/popular?api_key=22f697130e659aa6a1da22122fda7827')
         .then(response => response.json())
         .then(data => {
            const movies = data.results;
            let movieHTML = '';
            movies.forEach(movie => {
                const movieDiv = document.createElement('div');
                movieDiv.classList.add('movie-list-item');
                movieHTML += `
                   <div class="media-description">
                      <div class="description">  
                       <img src="https://image.tmdb.org/t/p/w500${movie.backdrop_path}" alt="${movie.title}" class="movie-list-item-img">
                       <div class="status-dot not-downloaded"></div>
                       <button class="download-btn" onclick="downloadMovie('${movie.title}')"><i class="fa-solid fa-download"></i></button>
                      </div> 
                      <p class="title">${movie.title}<p>
                   </div>
                `;
            });
            document.getElementById('popular-movies').innerHTML = movieHTML;
         });
}

function fetchPopularTVShows() {
    fetch('https://api.themoviedb.org/3/tv/popular?api_key=22f697130e659aa6a1da22122fda7827')
         .then(response => response.json())
         .then(data => {
            const tvShows = data.results;
            let tvHTML = '';
            tvShows.forEach(tv => {
                const tvDiv = document.createElement('div');
                tvDiv.classList.add('movie-list-item');
                tvHTML += `
                <div class="media-description">
                    <div class="description">
                       <img src="https://image.tmdb.org/t/p/w500${tv.poster_path}" alt="${tv.title}" class="movie-list-item-img">
                       <div class="status-dot not-downloaded"></div>
                       <button class="download-btn" onclick="downloadMovie('${tv.name}')"><i class="fa-solid fa-download"></i></button>
                   </div>
                   <p class="title">${tv.name}<p>
                 </div>  
                `;
            });
            document.getElementById('popular-tv').innerHTML = tvHTML;
         });
}

async function fetchPopularAnime(accessToken) {
    const query = `
    query {
  Page(page: 1, perPage: 10) { 
    media(type: ANIME, sort: POPULARITY_DESC) {
      title {
        romaji
        english
        native
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

        let animeHTML = '';
        animeList.forEach(anime => {
            const animeDiv = document.createElement('div');
            animeDiv.classList.add('movie-list-item');
            animeHTML += `
            <div class="media-description">
                    <div class="description">
                       <img src="${anime.coverImage.large}" alt="${anime.title.romanji}" class="movie-list-item-img">
                       <div class="status-dot not-downloaded"></div>
                       <button class="download-btn"onclick="setPopup()"><i class="fa-solid fa-download"></i></button>
                   </div>
                   <p class="title">${anime.title.english}<p>
            </div>
            `;
        });
        document.getElementById('popular-anime').innerHTML = animeHTML;
    } catch (error) {
        console.error('Error fetching popular anime:', error);
    }
}

/*.......................................................... Deals with dark theme ......................................................... */
const themeSwitch = document.querySelector('.switch input');
const htmlElement = document.documentElement;

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

// Toggle dark mode on checkbox change
themeSwitch.addEventListener('change', () => {
    if (themeSwitch.checked) {
        htmlElement.classList.add('dark-mode');
        localStorage.setItem('darkMode', 'true'); // Save preference
    } else {
        htmlElement.classList.remove('dark-mode');
        localStorage.setItem('darkMode', 'false'); // Save preference
    }
});

/*........................................................ Notification sidebar ............................................................... */
const notifSidebar = document.querySelector('.notification-sidebar');
const notifIcon = document.querySelector('.notif-icon');
const closeNotifButton = document.querySelector('.close-notif-btn');

// Toggle Notification Sidebar
notifIcon.addEventListener('click', () => {
    notifSidebar.classList.add('open');
});

// Close Notification Sidebar
closeNotifButton.addEventListener('click', () => {
    notifSidebar.classList.remove('open');
});

// Close Notification Sidebar on Outside Click
document.addEventListener('click', (event) => {
    if (
        !notifSidebar.contains(event.target) &&
        !notifIcon.contains(event.target)
    ) {
        notifSidebar.classList.remove('open');
    }
});

/*...................................................... Download Activity ............................................................. */
const downloadBtn = document.getElementById('download-btn');
const popupContainer = document.getElementById('popup-container');
const closePopupBtn = document.getElementById('close-popup');
const rssSelector = document.getElementById('rss-feed-selector');
const rssResults = document.getElementById('rss-results');

function setPopup() {
    popupContainer.classList.remove('hidden');
}

closePopupBtn.addEventListener('click', () => {
    popupContainer.classList.add('hidden');
});

// Fetch RSS feed data and display tabs
rssSelector.addEventListener('change', async () => {
    const feedName = rssSelector.value;
    const response = await fetch(`/rss/${feedName}`);
    const data = await response.json();

    rssResults.innerHTML = ''; // Clear previous results
    data.forEach(item => {
        const tab = document.createElement('div');
        tab.className = 'rss-tab';
        tab.textContent = item.title;
        tab.addEventListener('click', () => {
            window.open(item.magnetLink, '_blank');
        });
        rssResults.appendChild(tab);
    });
});

function downloadMovie(title) {
    alert(`Downloading movie: ${title}`)
}

function downloadTV(title) {
    alert(`Downloading TV Show: ${title}`)
}