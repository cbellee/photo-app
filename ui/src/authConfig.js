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
        clientId: "19348540-a082-4d89-9e2f-db0f9d3e7701",
        authority: "https://login.microsoftonline.com/kainiindustries.net",
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
    ]
};


export const storageRequestScope = {
    scopes: ["https://storage.azure.com/user_impersonation"]
};

export const apiConfig = {
    tenantId: "3d49be6f-6e38-404b-bbd4-f61c1a2d25bf",
    storageApiEndpoint: "https://storaodiwusgqeyiu.blob.core.windows.net",
    collectionApiEndpoint: "https://photo.orangecoast-f46e9cc8.australiaeast.azurecontainerapps.io/collections",
    albumApiEndpoint: "https://photo.orangecoast-f46e9cc8.australiaeast.azurecontainerapps.io/albums",
    photoApiEndpoint: "https://photo.orangecoast-f46e9cc8.australiaeast.azurecontainerapps.io/photos"
}
