// ---------------------------
// Types
// ---------------------------

import dayjs from 'dayjs';

export type EntityType = 'activity' | 'note' | 'friend';

export type FeedItem = {
    id: string;
    type: EntityType;
    desc: string;
    friends?: string[];
    tags?: string[];
    locations?: string[];
    createdAt: string; // ISO
};

// ---------------------------
// Mock data
// ---------------------------

export const MOCK_TAGS = ['#coworker', '#prank', '#wedding', '#family', '#office', '#fitness'];
export const MOCK_LOCATIONS = ['@Scranton', '@NYC', '@NiagaraFalls', '@Remote'];
export const MOCK_FRIENDS = ['Michael Scott', 'Jim Halpert', 'Pam Beesly', 'Dwight Schrute'];

export const INITIAL_FEED: FeedItem[] = [
    {
        id: '1',
        type: 'activity',
        desc: "Hosted Diversity Day workshop",
        friends: ['Michael Scott'],
        tags: ['#office'],
        locations: ['@Scranton'],
        createdAt: dayjs().subtract(1, 'day').toISOString(),
    },
    {
        id: '2',
        type: 'note',
        desc: "Put Dwight's stuff in Jell-O",
        friends: ['Jim Halpert'],
        tags: ['#prank'],
        locations: ['@Scranton'],
        createdAt: dayjs().subtract(2, 'day').toISOString(),
    },
    {
        id: '3',
        type: 'activity',
        desc: 'Engaged to Jim 💍',
        friends: ['Pam Beesly', 'Jim Halpert'],
        tags: ['#wedding'],
        locations: ['@NiagaraFalls'],
        createdAt: '2009-09-08T12:00:00.000Z',
    },
];
