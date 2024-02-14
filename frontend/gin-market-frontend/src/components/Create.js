import React, { useState } from 'react';

const Create = () => {
    const [itemName, setItemName] = useState('');
    const [itemPrice, setItemPrice] = useState('');
    const [itemDescription, setItemDescription] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();

        // TODO Signupしたユーザーのtokenを保存できるようにする
        console.log(localStorage.getItem('accessToken'))

        // Create機能: 新しいアイテムを作成し、JSON形式で返す
        fetch('http://localhost:8080/items', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('accessToken')}` // ユーザー認証が必要な場合
            },
            body: JSON.stringify({
                Name: itemName,
                Price: parseFloat(itemPrice), // 数値型に変換する
                Description: itemDescription
            })
        })
            .then(response => response.json())
            .then(data => {
                console.log('New item created:', data);
                // 成功した場合の処理をここに追加する
            })
            .catch(error => {
                console.error('Error creating item:', error);
                // エラーが発生した場合の処理をここに追加する
            });
    };

    return (
        <div>
            <h2>Create New Item</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Name:</label>
                    <input type="text" value={itemName} onChange={(e) => setItemName(e.target.value)} required />
                </div>
                <div>
                    <label>Price:</label>
                    <input type="number" value={itemPrice} onChange={(e) => setItemPrice(e.target.value)} required />
                </div>
                <div>
                    <label>Description:</label>
                    <input type="text" value={itemDescription} onChange={(e) => setItemDescription(e.target.value)} />
                </div>
                <button type="submit">Create</button>
            </form>
        </div>
    );
};

export default Create;
