import React, {useState} from "react";
import {
    Autocomplete,
    Button,
    Group,
    Kbd,
    MultiSelect,
    Paper,
    SegmentedControl,
    Stack,
    Text, Textarea
} from "@mantine/core";
import {useHotkeys, useLocalStorage} from "@mantine/hooks";
import {EntityType, FeedItem} from "@/mocks";

const MOCK_TAGS = ['#coworker', '#prank', '#wedding', '#family', '#office', '#fitness'];
const MOCK_LOCATIONS = ['@Scranton', '@NYC', '@NiagaraFalls', '@Remote'];
const MOCK_FRIENDS = ['Michael Scott', 'Jim Halpert', 'Pam Beesly', 'Dwight Schrute'];

export const QuickCapture = ({ onCreate }: { onCreate: (item: FeedItem) => void }) => {
    const [type, setType] = useLocalStorage<EntityType>({ key: 'frens.quick.type', defaultValue: 'note' });
    const [desc, setDesc] = useState('');
    const [friends, setFriends] = useState<string[]>([]);
    const [tags, setTags] = useState<string[]>([]);
    const [locations, setLocations] = useState<string[]>([]);

    const canSubmit = desc.trim().length > 0 || friends.length > 0 || tags.length > 0 || locations.length > 0;

    useHotkeys([
        ['f', () => setType('friend')],
        ['a', () => setType('activity')],
        ['n', () => setType('note')],
        ['mod+Enter', () => handleSubmit()],
    ]);

    const handleSubmit = () => {
        if (!canSubmit) return;
        if (type === 'friend') {
            // Minimal friend capture: use desc as friend name
            const name = desc.trim() || 'New Friend';
            onCreate({ id: crypto.randomUUID(), type: 'friend', desc: name, createdAt: new Date().toISOString() });
        } else {
            onCreate({
                id: crypto.randomUUID(),
                type,
                desc: desc.trim(),
                friends,
                tags,
                locations,
                createdAt: new Date().toISOString(),
            });
        }
        setDesc('');
        setFriends([]);
        setTags([]);
        setLocations([]);
    };

    return (
        <Paper withBorder p="md" radius="lg">
            <Group justify="space-between" align="flex-start" wrap="nowrap">
                <Stack flex={1} gap="sm">
                    <Group gap="xs" wrap="wrap">
                        <SegmentedControl
                            value={type}
                            onChange={(v) => setType(v as EntityType)}
                            data={[
                                { label: 'Friend', value: 'friend' },
                                { label: 'Activity', value: 'activity' },
                                { label: 'Note', value: 'note' },
                            ]}
                        />
                        <Text c="dimmed" size="sm">
                            Shortcuts: <Kbd>f</Kbd>/<Kbd>a</Kbd>/<Kbd>n</Kbd>/, submit <Kbd>⌘</Kbd>+<Kbd>Enter</Kbd>
                        </Text>
                    </Group>

                    {type !== 'friend' && (
                        <Group grow>
                            <Autocomplete
                                label="Friends"
                                placeholder="Type names…"
                                data={MOCK_FRIENDS}
                                onOptionSubmit={(value) => !friends.includes(value) && setFriends([...friends, value])}
                            />
                            <MultiSelect
                                label="Tags"
                                placeholder="#tag"
                                searchable
                                data={MOCK_TAGS}
                                value={tags}
                                onChange={setTags}
                            />
                            <MultiSelect
                                label="Locations"
                                placeholder="@location"
                                searchable
                                data={MOCK_LOCATIONS}
                                value={locations}
                                onChange={setLocations}
                            />
                        </Group>
                    )}

                    <Textarea
                        autosize
                        minRows={4}
                        maxRows={10}
                        placeholder={type === 'friend' ? 'Friend name (aliases later)…' : "What's on your mind?"}
                        value={desc}
                        onChange={(e) => setDesc(e.currentTarget.value)}
                    />

                    <Group justify="space-between">
                        <Text c="dimmed" size="sm">{type === 'friend' ? 'Creates a Friend profile' : 'Saved as Today by default'}</Text>
                        <Group>
                            <Button variant="subtle" onClick={() => { setDesc(''); setFriends([]); setTags([]); setLocations([]); }}>Clear</Button>
                            <Button onClick={handleSubmit} disabled={!canSubmit}>Save</Button>
                        </Group>
                    </Group>
                </Stack>
            </Group>
        </Paper>
    );
}
