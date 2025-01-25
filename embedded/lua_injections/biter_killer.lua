-- Function to handle the cleanup logic
local function cleanup_biters(player)
    local surface = player.surface
    local destroyed_count = 0

    -- Fun Star Trek Anecdotes
    local anecdotes = {
        "Spock would find this cleanup 'highly logical.'",
        "Captain Kirk just gave the order to 'exterminate all enemies!'",
        "Make it so, Number One! Biters will be no more.",
        "Worf suggests we show the biters 'honor'... by destroying them all.",
        "Scotty's working on the warp core, but he says, 'I cannae destroy biters any faster, Cap'n!'",
        "Resistance is futile, biters. You will be eliminated.",
        "Dr. McCoy reminds you, 'I’m a doctor, not a biter exterminator!'",
        "If Q were here, he’d snap his fingers and clean this up in no time.",
        "Biters aren’t part of the prime directive—engage!"
    }

    -- Iterate over all chunks on the surface
    for chunk in surface.get_chunks() do
        local area = {
            {chunk.x * 32, chunk.y * 32},
            {chunk.x * 32 + 32, chunk.y * 32 + 32}
        }

        -- Find all enemy entities (biters, spawners, worms) in the chunk
        local enemies = surface.find_entities_filtered({
            area = area,
            force = "enemy"
        })

        -- Destroy each enemy entity found
        for _, enemy in pairs(enemies) do
            if enemy and enemy.valid then
                game.print("[INFO] Destroying: " .. enemy.name .. " at (" .. math.floor(enemy.position.x) .. ", " .. math.floor(enemy.position.y) .. ")")
                enemy.destroy()
                destroyed_count = destroyed_count + 1
            end
        end
    end

    -- Select a random Star Trek anecdote
    local random_anecdote = anecdotes[math.random(1, #anecdotes)]

    -- Final report with humor
    game.print("[SUCCESS] Cleanup completed. Total enemies destroyed: " .. destroyed_count)
    game.print("[FUN] " .. random_anecdote)
end

-- Register a custom command to trigger the cleanup
commands.add_command("cleanup_biters", "Destroys all biters, spawners, and worms on the player's current surface.", function(cmd)
    local player = game.get_player(cmd.player_index)

    -- Ensure the command is run by a valid player
    if not player then
        game.print("[ERROR] Command can only be run by a player.")
        return
    end

    game.print("[INFO] Cleanup command triggered by " .. player.name .. ". Starting cleanup...")
    cleanup_biters(player)
end)