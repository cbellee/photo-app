import React, { useState, useEffect } from 'react'
import { getCollections } from '../photoService';
import { useMsal } from "@azure/msal-react";
import { photoApi } from "../authConfig";
import { Photo } from "./Photo"

export const CollectionGrid = () => {
    const { instance, accounts } = useMsal();
    const [collectionData, setCollectionData] = useState([])

    useEffect(() => {
        instance.acquireTokenSilent({
            ...photoApi,
            account: accounts[0]
        }).then((response) => {
            let accessToken = response.accessToken;
            getCollections(accessToken)
                .catch(error => console.log(error), (data = '[]') => console.log(data))
                .then((data) => setCollectionData(data));
            console.log("data: " + collectionData);
        })
    }, [])
    return (
        <div>
            {collectionData.map((collection) => (
                <Photo
                    name={collection.name}
                />
            ))}
        </div>
    );
}