import { Route, Routes } from 'react-router-dom';
import CollectionList from './CollectionList';
import AlbumList from './AlbumList';
import PhotoList from './PhotoList';
import { ImageUpload } from './ImageUpload';

function Pages() {
    return (
        <Routes>
            <Route path="/" element={<CollectionList />} />
            <Route path="/:collection" element={<AlbumList />} />
            <Route path="/:collection/:album" element={<PhotoList />} />
            <Route path="/uploads" element={<ImageUpload />} />
        </Routes>
    );
}

export default Pages;
