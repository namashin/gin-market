import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

const FindById = () => {
    const { id } = useParams(); // URLパラメーターのidを取得
    const [item, setItem] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchItemById = async () => {
            try {
                const accessToken = localStorage.getItem('accessToken');
                const response = await fetch(`http://localhost:8080/items/${id}`, {
                    headers: {
                        'Authorization': `Bearer ${accessToken}`
                    }
                });

                // responseをjsonにして表示
                const data = await response.json();
                console.log(data["data"]);

                setItem(data["data"])

            } catch (error) {
                console.error('Error fetching item:', error);
                setError(error.toString());
            }
        };

        fetchItemById();
    }, [id]);

    return (
        <div>
            <h2>My Items</h2>
            <pre>{JSON.stringify(item)}</pre>
        </div>
    );
};

export default FindById;
