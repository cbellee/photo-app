/*
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License.
 */

import { LogLevel } from "@azure/msal-browser";

/**
 * Configuration object to be passed to MSAL instance on creation. 
 * For a full list of MSAL.js configuration parameters, visit:
 * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-browser/docs/configuration.md 
 */
export const msalConfig = {
    auth: {
        clientId: "{{CLIENT_ID}}",
        authority: "https://login.microsoftonline.com/{{DOMAIN_NAME}}",
        redirectUri: "http://localhost:3000"
    },
    cache: {
        cacheLocation: "sessionStorage", // This configures where your cache will be stored
        storeAuthStateInCookie: false, // Set this to "true" if you are having issues on IE11 or Edge
    },
    system: {
        loggerOptions: {
            loggerCallback: (level, message, containsPii) => {
                if (containsPii) {
                    return;
                }
                switch (level) {
                    case LogLevel.Error:
                        console.error(message);
                        return;
                    case LogLevel.Info:
                        console.info(message);
                        return;
                    case LogLevel.Verbose:
                        console.debug(message);
                        return;
                    case LogLevel.Warning:
                        console.warn(message);
                        return;
                }
            }
        }
    }
};

/**
 * Scopes you add here will be prompted for user consent during sign-in.
 * By default, MSAL.js will add OIDC scopes (openid, profile, email) to any login request.
 * For more information about OIDC scopes, visit: 
 * https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-permissions-and-consent#openid-connect-scopes
 */
export const loginRequest = {
    scopes: [
        "{{PHOTO_READ_SCOPE}}",
        "{{PHOTO_WRITE_SCOPE}}"
    ]
};

export const uploadApi = {
    scopes: [
        "{{UPLOAD_READ_SCOPE}}",
        "{{UPLOAD_WRITE_SCOPE}}"
    ],
}

export const photoApi = {
    scopes: [
        "{{PHOTO_READ_SCOPE}}",
        "{{PHOTO_WRITE_SCOPE}}"
    ],
}

export const apiConfig = {
    uploadApiEndpoint: "{{UPLOAD_API_ENDPOINT}}",
    collectionApiEndpoint: "{{COLLECTION_API_ENDPOINT}}",
    albumApiEndpoint: "{{ALBUM_API_ENDPOINT}}",
    photoApiEndpoint: "{{PHOTO_API_ENDPOINT}}"
}
