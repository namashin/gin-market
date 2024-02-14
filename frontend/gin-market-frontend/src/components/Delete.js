import React from 'react';
import { useParams } from 'react-router-dom';

const Delete = () => {
    const { id } = useParams(); // URLパラメーターのidを取得

    const handleDelete = () => {
        const accessToken = localStorage.getItem('accessToken');

        fetch(`http://localhost:8080/items/${id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}`
            }
        })
            .then(response => {
                if (response.ok) {
                    console.log('Item deleted successfully');
                } else {
                    throw new Error('Failed to delete item');
                }
            })
            .catch(error => {
                console.error('Error deleting item:', error);
            });
    };

    return (
        <div>
            <button onClick={handleDelete}>Delete</button>
        </div>
    );
};

export default Delete;
