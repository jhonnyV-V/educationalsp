local client = vim.lsp.start_client({
	name = "educationalsp",
	cmd = { "example/pat/main" },
	-- on_attach = on_attach,
})
if not client then
	vim.notify("Failed to start educationalsp", vim.log.levels.ERROR)
	return
end

vim.api.nvim_create_autocmd("FileType", {
	pattern = "markdown",
	callback = function()
		vim.lsp.buf_attach_client(0, client)
	end,
})
