// API client for communicating with the backend

const API_BASE = "/api";

export interface Contact {
  id: string;
  type: string;
  value: string;
  tags?: string[];
}

export interface Friend {
  id: string;
  name: string;
  description?: string;
  nicknames?: string[];
  tags?: string[];
  locations?: string[];
  contacts?: Contact[];
  activitiesCount: number;
  notesCount: number;
  lastActivity?: string;
}

export interface Stats {
  friends: number;
  locations: number;
  activities: number;
  notes: number;
}

export interface RankedItem {
  id: string;
  name: string;
  count: number;
  lastActivity?: string;
}

export interface TimelineDataPoint {
  month: string;
  activities: number;
  notes: number;
}

export interface Insight {
  type: "reconnect" | "streak" | "milestone";
  title: string;
  description: string;
  friendId?: string;
}

export interface ComprehensiveStats {
  counts: Stats;
  topFriends: RankedItem[];
  topLocations: RankedItem[];
  topTags: RankedItem[];
  activityTimeline: TimelineDataPoint[];
  insights: Insight[];
}

export interface SyncStatus {
  gitInstalled: boolean;
  gitInited: boolean;
  branch?: string;
  hasChanges: boolean;
  changeCount: number;
}

export type FeedItemType = "activity" | "note" | "friend_added" | "location_added";

export interface FeedItem {
  id: string;
  type: FeedItemType;
  date: string;
  description: string;
  friendIds?: string[];
  locationIds?: string[];
  tags?: string[];
  entityId?: string;
  entityName?: string;
}

export type EventType = "activity" | "note";

export interface Event {
  id: string;
  type: EventType;
  date: string;
  description: string;
  friendIds?: string[];
  locationIds?: string[];
  tags?: string[];
}

export interface Location {
  id: string;
  name: string;
  country?: string;
  description?: string;
  aliases?: string[];
  tags?: string[];
  activitiesCount: number;
  notesCount: number;
  lastActivity?: string;
}

class ApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function fetchJson<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`);

  if (!response.ok) {
    const text = await response.text();
    throw new ApiError(response.status, text || response.statusText);
  }

  return response.json();
}

export const api = {
  friends: {
    list: (): Promise<Friend[]> => fetchJson<Friend[]>("/friends"),
    get: (id: string): Promise<Friend> => fetchJson<Friend>(`/friends/${id}`),
    activities: (id: string): Promise<Event[]> =>
      fetchJson<Event[]>(`/friends/${id}/activities`),
    notes: (id: string): Promise<Event[]> =>
      fetchJson<Event[]>(`/friends/${id}/notes`),
  },
  locations: {
    list: (): Promise<Location[]> => fetchJson<Location[]>("/locations"),
    get: (id: string): Promise<Location> => fetchJson<Location>(`/locations/${id}`),
    activities: (id: string): Promise<Event[]> =>
      fetchJson<Event[]>(`/locations/${id}/activities`),
    notes: (id: string): Promise<Event[]> =>
      fetchJson<Event[]>(`/locations/${id}/notes`),
  },
  notes: {
    list: (): Promise<Event[]> => fetchJson<Event[]>("/notes"),
  },
  activities: {
    list: (): Promise<Event[]> => fetchJson<Event[]>("/activities"),
  },
  stats: {
    get: (): Promise<Stats> => fetchJson<Stats>("/stats"),
    getComprehensive: (): Promise<ComprehensiveStats> =>
      fetchJson<ComprehensiveStats>("/stats/comprehensive"),
  },
  sync: {
    status: (): Promise<SyncStatus> => fetchJson<SyncStatus>("/sync/status"),
  },
  feed: {
    list: (): Promise<FeedItem[]> => fetchJson<FeedItem[]>("/feed"),
  },
};
