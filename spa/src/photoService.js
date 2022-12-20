import { apiConfig } from "./authConfig";

/**
 * Attaches a given access token to a MS Graph API call. Returns information about the user
 * @param accessToken 
 */

export async function getCollections(accessToken) {
    const headers = new Headers();
    const bearer = `Bearer ${accessToken}`;

    headers.append("Authorization", bearer);
    headers.append("Access-Control-Allow-Headers", "*");
    headers.append("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE");
    headers.append("Access-Control-Allow-Origin", "http://localhost:3000");

    const options = {
        method: "GET",
        headers: headers,
    };

    return fetch(apiConfig.collectionApiEndpoint, options)
        .then(response => response.json())
        .catch(error => console.log(error));
}

export async function getAlbums(accessToken) {
    const headers = new Headers();
    const bearer = `Bearer ${accessToken}`;

    headers.append("Authorization", bearer);
    headers.append("Access-Control-Allow-Headers", "*");
    headers.append("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE");
    headers.append("Access-Control-Allow-Origin", "http://localhost:3000");

    const options = {
        method: "GET",
        headers: headers,
    };

    return fetch(apiConfig.albumApiEndpoint, options)
        .then(response => response.json())
        .catch(error => console.log(error));
}

export async function getAlbumPhotos(accessToken) {
    const headers = new Headers();
    const bearer = `Bearer ${accessToken}`;

    headers.append("Authorization", bearer);
    headers.append("Access-Control-Allow-Headers", "*");
    headers.append("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE");
    headers.append("Access-Control-Allow-Origin", "http://localhost:3000");

    const options = {
        method: "GET",
        headers: headers,
    };

    return fetch(apiConfig.photoApiEndpoint + "?collection=2022&album=sport", options)
        .then(response => response.json())
        .catch(error => console.log(error));
}
