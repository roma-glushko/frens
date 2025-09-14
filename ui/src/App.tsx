import React, { useMemo, useState } from 'react';
import {
    AppShell,
    AppShellHeader,
    AppShellNavbar,
    AppShellAside,
    AppShellMain,
    Burger,
    Group,
    Title,
    Text,
    Button,
    ActionIcon,
    Divider,
    ScrollArea,
    Menu,
    Tooltip,
    TextInput,
    MultiSelect,
    Paper,
    Stack,
    Badge,
    Avatar,
    useMantineColorScheme,
    MantineProvider,
    ColorSchemeScript,
    SegmentedControl,
} from '@mantine/core';
import { useDisclosure, useHotkeys } from '@mantine/hooks';
import { IconPlus, IconDots, IconMoon, IconSun, IconSettings, IconSearch, IconUsers, IconMapPin, IconNotebook, IconTag, IconBell, IconHome2, IconCalendar, IconFilter } from '@tabler/icons-react';

import {QuickCapture} from "@/components/QuickCapture";
import {EntityType, FeedItem, INITIAL_FEED, MOCK_FRIENDS, MOCK_LOCATIONS, MOCK_TAGS} from "@/mocks";
import dayjs from "dayjs";
// ---------------------------
// Header elements
// ---------------------------

function ThemeToggle() {
    const { colorScheme, setColorScheme } = useMantineColorScheme();
    const dark = colorScheme === 'dark';
    return (
        <Tooltip label={dark ? 'Switch to light' : 'Switch to dark'}>
            <ActionIcon
                variant="subtle"
                onClick={() => setColorScheme(dark ? 'light' : 'dark')}
                aria-label="Toggle color scheme"
                size="lg"
            >
                {dark ? <IconSun /> : <IconMoon />}
            </ActionIcon>
        </Tooltip>
    );
}

function GlobalSearch({ onSubmit }: { onSubmit: (q: string) => void }) {
    const [q, setQ] = useState('');
    useHotkeys([
        ['/', (e) => { e.preventDefault(); const el = document.getElementById('global-search'); el?.focus(); }],
        ['Enter', () => { const el = document.activeElement as HTMLElement | null; if (el?.id === 'global-search') onSubmit(q); }],
    ]);
    return (
        <TextInput
            id="global-search"
            leftSection={<IconSearch size={16} />}
            placeholder="Search friends, tags, locations… (/ to focus)"
            value={q}
            onChange={(e) => setQ(e.currentTarget.value)}
            w={420}
            radius="md"
        />
    );
}

function AddMenu({ onAdd }: { onAdd: (type: EntityType) => void }) {
    return (
        <Menu shadow="md" width={200}>
            <Menu.Target>
                <Button leftSection={<IconPlus size={16} />} radius="md">Add</Button>
            </Menu.Target>
            <Menu.Dropdown>
                <Menu.Item leftSection={<IconNotebook size={16} />} onClick={() => onAdd('activity')}>Activity</Menu.Item>
                <Menu.Item leftSection={<IconNotebook size={16} />} onClick={() => onAdd('note')}>Note</Menu.Item>
                <Menu.Item leftSection={<IconUsers size={16} />} onClick={() => onAdd('friend')}>Friend</Menu.Item>
            </Menu.Dropdown>
        </Menu>
    );
}

// ---------------------------
// Sidebar nav
// ---------------------------

function NavLink({ icon, label, active = false }: { icon: React.ReactNode; label: string; active?: boolean }) {
    return (
        <Button
            variant={active ? 'light' : 'subtle'}
            leftSection={icon}
            fullWidth
            justify="flex-start"
            radius="md"
        >
            {label}
        </Button>
    );
}

function Sidebar() {
    return (
        <Stack p="md" gap="xs">
            <NavLink icon={<IconHome2 size={18} />} label="Home" active />
            <NavLink icon={<IconUsers size={18} />} label="Friends" />
            <NavLink icon={<IconMapPin size={18} />} label="Locations" />
            <NavLink icon={<IconNotebook size={18} />} label="Activities" />
            <NavLink icon={<IconNotebook size={18} />} label="Notes" />
            <NavLink icon={<IconTag size={18} />} label="Tags" />
            <NavLink icon={<IconCalendar size={18} />} label="Calendar" />
            <NavLink icon={<IconBell size={18} />} label="Reminders" />
        </Stack>
    );
}

// ---------------------------
// Filters (Aside)
// ---------------------------

function FiltersAside({
                          filters,
                          setFilters,
                      }: {
    filters: { friends: string[]; tags: string[]; locations: string[]; onlyNotes: boolean; onlyActivities: boolean };
    setFilters: (f: FiltersAsideProps['filters']) => void;
}) {
    const update = (patch: Partial<FiltersAsideProps['filters']>) => setFilters({ ...filters, ...patch });
    return (
        <ScrollArea p="md" h="100%">
            <Group justify="space-between" mb="xs">
                <Text fw={600}>Filters</Text>
                <IconFilter size={16} />
            </Group>
            <Divider mb="sm" />
            <MultiSelect
                label="Friends"
                data={MOCK_FRIENDS}
                searchable
                value={filters.friends}
                onChange={(v) => update({ friends: v })}
                placeholder="Pick friends"
                mb="sm"
            />
            <MultiSelect
                label="Tags"
                data={MOCK_TAGS}
                searchable
                value={filters.tags}
                onChange={(v) => update({ tags: v })}
                placeholder="Pick tags"
                mb="sm"
            />
            <MultiSelect
                label="Locations"
                data={MOCK_LOCATIONS}
                searchable
                value={filters.locations}
                onChange={(v) => update({ locations: v })}
                placeholder="Pick locations"
                mb="md"
            />
            <SegmentedControl
                fullWidth
                mb="sm"
                value={filters.onlyActivities ? 'activities' : filters.onlyNotes ? 'notes' : 'all'}
                onChange={(val) =>
                    update({
                        onlyActivities: val === 'activities',
                        onlyNotes: val === 'notes',
                    })
                }
                data={[
                    { label: 'All', value: 'all' },
                    { label: 'Activities', value: 'activities' },
                    { label: 'Notes', value: 'notes' },
                ]}
            />
            <Button variant="subtle" onClick={() => setFilters({ friends: [], tags: [], locations: [], onlyActivities: false, onlyNotes: false })}>
                Reset filters
            </Button>
        </ScrollArea>
    );
}

type FiltersAsideProps = React.ComponentProps<typeof FiltersAside>;

// ---------------------------
// Feed (Timeline)
// ---------------------------

function Feed({ items }: { items: FeedItem[] }) {
    const grouped = useMemo(() => {
        const byDay = new Map<string, FeedItem[]>();
        for (const it of items) {
            const key = dayjs(it.createdAt).format('YYYY-MM-DD');
            const arr = byDay.get(key) ?? [];
            arr.push(it);
            byDay.set(key, arr);
        }
        // Sort days desc
        const days = [...byDay.entries()].sort((a, b) => (a[0] < b[0] ? 1 : -1));
        // Sort items within day desc by time
        return days.map(([day, list]) => [day, list.sort((a, b) => (a.createdAt < b.createdAt ? 1 : -1))] as const);
    }, [items]);

    if (items.length === 0) return <Text c="dimmed">No entries yet. Start by writing something above.</Text>;

    return (
        <Stack>
            {grouped.map(([day, list]) => (
                <Stack key={day} gap="xs">
                    <Group>
                        <Text fw={700}>{dayjs(day).calendar?.() ?? dayjs(day).format('MMM D, YYYY')}</Text>
                        <Divider style={{ flex: 1 }} />
                    </Group>
                    {list.map((it) => (
                        <Paper key={it.id} withBorder p="md" radius="md">
                            <Group align="flex-start" gap="sm">
                                <Avatar radius="xl">
                                    {it.type === 'friend' ? '👤' : it.type === 'activity' ? '⚡' : '📝'}
                                </Avatar>
                                <Stack gap={4} flex={1}>
                                    <Group wrap="wrap" gap={6}>
                                        <Text fw={600}>{it.type.toUpperCase()}</Text>
                                        <Text c="dimmed">{dayjs(it.createdAt).format('HH:mm')}</Text>
                                    </Group>
                                    <Text>{it.desc}</Text>
                                    <Group gap={6} wrap="wrap">
                                        {it.friends?.map((f) => (
                                            <Badge key={f} variant="light" leftSection={<IconUsers size={14} />}>{f}</Badge>
                                        ))}
                                        {it.tags?.map((t) => (
                                            <Badge key={t} variant="outline">{t}</Badge>
                                        ))}
                                        {it.locations?.map((l) => (
                                            <Badge key={l} variant="light" leftSection={<IconMapPin size={14} />}>{l}</Badge>
                                        ))}
                                    </Group>
                                </Stack>
                                <Menu withinPortal position="bottom-end">
                                    <Menu.Target>
                                        <ActionIcon variant="subtle"><IconDots /></ActionIcon>
                                    </Menu.Target>
                                    <Menu.Dropdown>
                                        <Menu.Item>Edit</Menu.Item>
                                        <Menu.Item c="red">Delete</Menu.Item>
                                    </Menu.Dropdown>
                                </Menu>
                            </Group>
                        </Paper>
                    ))}
                </Stack>
            ))}
        </Stack>
    );
}

// ---------------------------
// Root App
// ---------------------------

function Shell() {
    const [opened, { toggle }] = useDisclosure(true);
    const [asideOpened, setAsideOpened] = useState(true);
    const [feed, setFeed] = useState<FeedItem[]>(INITIAL_FEED);
    const [filters, setFilters] = useState<FiltersAsideProps['filters']>({ friends: [], tags: [], locations: [], onlyActivities: false, onlyNotes: false });

    const filteredFeed = useMemo(() => {
        return feed.filter((it) => {
            if (filters.onlyActivities && it.type !== 'activity') return false;
            if (filters.onlyNotes && it.type !== 'note') return false;
            if (filters.friends.length && !(it.friends ?? []).some((f) => filters.friends.includes(f))) return false;
            if (filters.tags.length && !(it.tags ?? []).some((t) => filters.tags.includes(t))) return false;
            if (filters.locations.length && !(it.locations ?? []).some((l) => filters.locations.includes(l))) return false;
            return true;
        });
    }, [feed, filters]);

    const addItem = (item: FeedItem) => setFeed((prev) => [item, ...prev]);

    const handleSearch = (q: string) => {
        // naive search over description & tags
        const matches = feed.filter((it) =>
            it.desc.toLowerCase().includes(q.toLowerCase()) || (it.tags ?? []).some((t) => t.toLowerCase().includes(q.toLowerCase()))
        );
        // For a real app, route to a /search page. Here we just replace list temporarily.
        // no-op — or you could set a local state to show matches.
        console.log('Search matches', matches.length);
    };

    return (
        <AppShell
            header={{ height: 60 }}
            navbar={{ width: 240, breakpoint: 'sm', collapsed: { mobile: !opened } }}
            aside={{ width: 300, breakpoint: 'lg', collapsed: { desktop: !asideOpened } }}
            padding="md"
        >
            <AppShellHeader>
                <Group h="100%" px="md" justify="space-between">
                    <Group>
                        <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
                        <Title order={4}>Frens</Title>
                    </Group>
                    <GlobalSearch onSubmit={handleSearch} />
                    <Group>
                        <AddMenu onAdd={() => { /* focus composer? */ }} />
                        <ThemeToggle />
                        <ActionIcon variant="subtle" aria-label="Settings"><IconSettings /></ActionIcon>
                    </Group>
                </Group>
            </AppShellHeader>

            <AppShellNavbar>
                <ScrollArea h="100%"><Sidebar /></ScrollArea>
            </AppShellNavbar>

            <AppShellAside>
                <FiltersAside filters={filters} setFilters={setFilters} />
            </AppShellAside>

            <AppShellMain>
                <Stack gap="md">
                    <QuickCapture onCreate={addItem} />
                    <Feed items={filteredFeed} />
                </Stack>
            </AppShellMain>
        </AppShell>
    );
}

export default function FrensApp() {
    return (
        <MantineProvider defaultColorScheme="auto">
            <ColorSchemeScript />
            <Shell />
        </MantineProvider>
    );
}
