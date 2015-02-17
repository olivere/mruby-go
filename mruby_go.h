// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

#ifndef MRUBY_GO_H
#define MRUBY_GO_H

#include <stdlib.h>
#include <string.h>

#include <mruby.h>
#include <mruby/array.h>
#include <mruby/class.h>
#include <mruby/hash.h>
#include <mruby/proc.h>
#include <mruby/data.h>
#include <mruby/compile.h>
#include <mruby/string.h>
#include <mruby/value.h>
#include <mruby/variable.h>

static inline struct mrbc_context *my_context_new(mrb_state *mrb, const char *filename, mrb_bool capture_errors, mrb_bool no_exec) {
	mrbc_context *ctx;

  ctx = mrbc_context_new(mrb);
  ctx->capture_errors = capture_errors;
  ctx->lineno = 1;
  ctx->no_exec = no_exec;
  mrbc_filename(mrb, ctx, filename);

	return ctx;
}

static inline struct mrb_parser_state *my_parse(mrb_state *mrb, mrbc_context *ctx, char *ruby_code) {
	struct mrb_parser_state *parser = mrb_parser_new(mrb);

	parser->s = ruby_code;
	parser->send = ruby_code + strlen(ruby_code);
	parser->lineno = 1;
	mrb_parser_parse(parser, ctx);

	return parser;
}

//mrb_value
//my_run(mrb_state *mrb, int n) {
//	return mrb_run(mrb,
//		mrb_proc_new(mrb, mrb->irep[n]),
//		mrb_top_self(mrb)
//	);
//}

static inline mrb_value my_run(mrb_state *mrb, struct RProc *proc) {
	return mrb_context_run(mrb,
		proc,
		mrb_top_self(mrb),
		proc->body.irep->nregs
	);
}

static inline int has_exception(mrb_state *mrb) {
	return mrb->exc != 0;
}

static inline void reset_exception(mrb_state *mrb) {
	mrb->exc = 0;
}

static inline const char *get_exception_message(mrb_state *mrb) {
	mrb_value val = mrb_obj_value(mrb->exc);
	return mrb_string_value_ptr(mrb, val);
}

// Value helpers

static inline int my_type(mrb_value v) {
	return mrb_type(v);
}

static inline mrb_value any_to_str(mrb_state *mrb, mrb_value v) {
	return mrb_any_to_s(mrb, v);
}

static inline mrb_float get_float(mrb_value v) {
	return mrb_float(v);
}

static inline mrb_value get_float_value(mrb_state *mrb, mrb_float f) {
	return mrb_float_value(mrb, f);
}

static inline mrb_int get_fixnum(mrb_value v) {
	return mrb_fixnum(v);
}

static inline mrb_sym get_symbol(mrb_value v) {
	return mrb_symbol(v);
}

static inline int is_fixnum(mrb_value v) {
	return mrb_fixnum_p(v);
}

static inline int is_float(mrb_value v) {
	return mrb_float_p(v);
}

static inline int is_undef(mrb_value v) {
	return mrb_undef_p(v);
}

static inline int is_nil(mrb_value v) {
	return mrb_nil_p(v);
}

static inline int is_symbol(mrb_value v) {
	return mrb_symbol_p(v);
}

static inline int is_array(mrb_value v) {
	return mrb_array_p(v);
}

static inline int is_string(mrb_value v) {
	return mrb_string_p(v);
}

static inline int is_hash(mrb_value v) {
	return mrb_hash_p(v);
}

static inline int is_cptr_p(mrb_value v) {
	return mrb_cptr_p(v);
}

static inline int is_bool(mrb_value v) {
	return mrb_bool(v) != 0;
}

static inline mrb_value get_ary_entry(mrb_value ary, int index) {
	return mrb_ary_entry(ary, ((mrb_int)index));
}

static inline struct RProc *my_mrb_proc_ptr(mrb_value v) {
	return mrb_proc_ptr(v);
}

static inline mrb_bool my_mrb_const_defined_at(mrb_state *mrb, const char *name) {
	mrb_sym sym = mrb_intern_cstr(mrb, name);
	return mrb_const_defined_at(mrb, mrb_obj_value(mrb->object_class), sym);
}

// Ruby -> Go

extern mrb_value my_mrb_func_call(mrb_state *, mrb_value);

static inline mrb_func_t my_mrb_func_call_t() {
	return &my_mrb_func_call;
}

/*
extern mrb_value my_mrb_class_func_call(mrb_state *, mrb_value);

static inline mrb_func_t my_mrb_class_func_call_t() {
	return &my_mrb_class_func_call;
}
*/

// Args

// required arguments
static inline mrb_aspec args_any() {
	return MRB_ARGS_ANY();
}

// required arguments
static inline mrb_aspec args_none() {
	return MRB_ARGS_NONE();
}

// required arguments
static inline mrb_aspec args_req(int n) {
	return MRB_ARGS_REQ(n);
}

// optional arguments
static inline mrb_aspec args_opt(int n) {
	return MRB_ARGS_OPT(n);
}

// required arguments
static inline mrb_aspec args_arg(int req, int opt) {
	return MRB_ARGS_ARG(req, opt);
}

static inline mrb_value my_get_args(mrb_state *mrb, mrb_value self, const char *format) {
	mrb_value arg;
	mrb_get_args(mrb, format, &arg);
	return arg;
}

#endif
