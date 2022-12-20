import React, { useState, useEffect } from 'react'
import { getAlbumPhotos } from '../photoService';
import { useMsal } from "@azure/msal-react";
import { Photo } from "./Photo.js"

export const PhotoGrid = () => {
    const { instance, accounts } = useMsal();
    const [photoData, setPhotoData] = useState([])

    useEffect(() => {
            getAlbumPhotos()
                .catch(error => console.log(error), (data = '[]') => console.log(data))
                .then((data) => setPhotoData(data));        
    }, [])
    return (
        <div>
            {photoData.map((photo) => (
                <Photo
                    name={photo.value.name}
                    url={photo.value.thumbsUrl}
                    id={photo.value.id}
                    key={photo.value.id}
                />
            ))}
        </div>
    );
}